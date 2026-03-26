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
	ON CONFLICT (size) DO NOTHING
	RETURNING id
`

func (r *Repo) CreatePack(ctx context.Context, pack *app.Pack) error {
	return r.db.QueryRowContext(ctx, CreatePackQuery, pack.ID, pack.Size).Scan(&pack.ID)
}
