package pg

import (
	"context"
	"database/sql"

	"github.com/hferr/pack-api/internal/app"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

const ListPacksQuery = `SELECT * FROM packs`

func (r *Repo) ListPacks(ctx context.Context) (app.Packs, error) {
	var packs app.Packs

	rows, err := r.db.QueryContext(ctx, ListPacksQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p app.Pack
		err := rows.Scan(
			&p.ID,
			&p.Size,
		)
		if err != nil {
			return nil, err
		}
		packs = append(packs, p)
	}

	return packs, nil
}

const CreatePackQuery = `
	INSERT INTO packs
		(id, size)
	VALUES
		($1, $2)
	ON CONFLICT (size) DO UPDATE SET size = EXCLUDED.size
	RETURNING id
`

func (r *Repo) CreatePack(ctx context.Context, pack *app.Pack) error {
	return r.db.QueryRowContext(ctx, CreatePackQuery, pack.ID, pack.Size).Scan(&pack.ID)
}

func (r *Repo) RebuildPacks(ctx context.Context, packs app.Packs) error {
	err := r.RunInTx(ctx, func(tx *sql.Tx) error {
		if err := r.deleteAllPacks(ctx); err != nil {
			return err
		}

		for _, p := range packs {
			if err := r.CreatePack(ctx, &p); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

const DeleteAllPacksQuery = `DELETE FROM packs`

func (r *Repo) deleteAllPacks(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, DeleteAllPacksQuery)
	return err
}

func (r *Repo) RunInTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
