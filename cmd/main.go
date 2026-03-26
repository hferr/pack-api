package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hferr/pack-api/config"
	"github.com/hferr/pack-api/internal/app"
	"github.com/hferr/pack-api/internal/httpjson"
	"github.com/hferr/pack-api/internal/repositories/pg"
	"github.com/hferr/pack-api/migrations"
)

const fmtDbConnString = "host=%s user=%s password=%s dbname=%s port=%d"

func main() {
	cfg := config.New()

	db, err := initPostgresDb(cfg)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer db.Close()

	repo := pg.NewRepo(db)

	packService := app.NewPackService(repo)

	h := httpjson.NewHandler(
		packService,
	)

	s := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:      h.NewRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func initPostgresDb(cfg *config.Cfg) (*sql.DB, error) {
	connString := fmt.Sprintf(
		fmtDbConnString,
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort,
	)

	psql, err := pg.NewPostgresDb(connString)
	if err != nil {
		return nil, err
	}

	if err := migrations.ApplyMigrations(psql.Db); err != nil {
		return nil, err
	}

	return psql.Db, nil
}
