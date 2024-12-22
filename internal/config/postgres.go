package config

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"time"
)

func NewDB(config *koanf.Koanf, log *zerolog.Logger) *sql.DB {
	dbUri := config.String("POSTGRES_URL")

	db, err := sql.Open("pgx", dbUri)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open postgresql")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping to postgresql")
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
