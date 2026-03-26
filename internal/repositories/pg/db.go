package pg

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	Db *sql.DB
}

func NewPostgresDb(connString string) (*Postgres, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection: %w", err)
	}

	return &Postgres{
		Db: db,
	}, nil
}
