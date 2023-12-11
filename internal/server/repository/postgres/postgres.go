package postgres

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// Config структура с полями для подключения к базе данных.
type Config struct {
	// строка подключения с базой данных.
	ConnURL string
	// максимальное количество открытых соединений с базой данных.
	MaxConn int
	// максимальное количество времени, в течение которого соединение может быть использовано повторно.
	MaxConnLifeTime time.Duration
	// максимальное количество времени, в течение которого соединение может простаивать.
	MaxConnIdleTime time.Duration
	// логгер.
	Logger zerolog.Logger
}

type Postgres struct {
	*sqlx.DB
}

func New(cfg Config) (Postgres, error) {
	l := cfg.Logger.With().Str("postgres", "New").Logger()

	db, err := sqlx.Connect("pgx", cfg.ConnURL)
	if err != nil {
		return Postgres{}, fmt.Errorf("postgres connection error: %w", err)
	}
	l.Info().Msg("succeeded in connecting to postgres")

	db.SetConnMaxLifetime(cfg.MaxConnLifeTime)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)
	db.SetMaxOpenConns(cfg.MaxConn)

	return Postgres{DB: db}, nil
}
