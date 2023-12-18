package postgres

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/1Asi1/metric-track.git/internal/server/models"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const (
	countStep      = 2
	retryStopCount = 5
)

var (
	count = 1
)

type Name string

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

//go:embed migrations/*.sql
var migrationsDir embed.FS

type storage interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	Ping() error
	Close() error
}

type Store struct {
	storage
	log zerolog.Logger
}

func New(cfg Config, log zerolog.Logger) (*Store, error) {
	l := cfg.Logger.With().Str("postgres", "New").Logger()
	var db *sqlx.DB
	var err error
	for ; ; count += countStep {
		ticker := time.NewTicker(time.Duration(count) * time.Second)
		db, err = sqlx.Connect("pgx", cfg.ConnURL)
		if err != nil {
			if _, ok := (err).(pgx.PgError); !ok {
				return &Store{}, fmt.Errorf("postgres connection error: %w", err)
			}

			pgErrCode := (err).(pgx.PgError).Code
			if pgErrCode == pgerrcode.InvalidAuthorizationSpecification {
				l.Info().Msgf("try connection sec: %d", count)
				<-ticker.C
				l.Err(err).Msg("sqlx.Connect try agan...")
				if count == retryStopCount {
					l.Error().Msg("sqlx.Connect try cancel")
					return &Store{}, fmt.Errorf("postgres connection error: %w", err)
				}
				continue
			}

			return &Store{}, fmt.Errorf("postgres connection error: %w", err)
		}
		break
	}

	l.Info().Msg("succeeded in connecting to postgres")

	db.SetConnMaxLifetime(cfg.MaxConnLifeTime)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)
	db.SetMaxOpenConns(cfg.MaxConn)

	if err = runMigrations(cfg.ConnURL); err != nil {
		return &Store{}, fmt.Errorf("runMigrations error: %w", err)
	}

	return &Store{storage: db, log: log}, nil
}

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
	}
	return nil
}

func (s *Store) Get(ctx context.Context) (map[string]memory.Type, error) {
	query := `
	SELECT
	    id,
	    gauge,
		counter
	FROM tbl_metrics
`
	var models []models.Metric
	if err := s.SelectContext(ctx, &models, query); err != nil {
		return nil, fmt.Errorf("Get %w;", err)
	}

	result := make(map[string]memory.Type)
	if len(models) != 0 {
		for _, v := range models {
			result[v.ID] = memory.Type{
				Gauge:   v.Gauge,
				Counter: v.Counter,
			}
		}
	}
	return result, nil
}

func (s *Store) GetOne(ctx context.Context, name string) (memory.Type, error) {
	query := `
	SELECT
	    gauge,
		counter
	FROM tbl_metrics
	WHERE id = $1
`
	var model memory.Type
	err := s.GetContext(ctx, &model, query, name)
	if err != nil {
		return memory.Type{}, fmt.Errorf("GetOne: %w", err)
	}

	return model, nil
}

func (s *Store) Update(ctx context.Context, data map[string]memory.Type) {

	l := s.log.With().Str("postgres", "Update").Logger()

	name := fmt.Sprintf("%v", ctx.Value(Name("name")))
	model := models.Metric{
		ID:      name,
		Gauge:   data[name].Gauge,
		Counter: data[name].Counter,
	}

	query := `
		INSERT INTO tbl_metrics(id,gauge,counter)
		VALUES (:id, :gauge, :counter)
		ON CONFLICT (id) DO UPDATE
		SET
		    gauge = EXCLUDED.gauge,
		    counter = EXCLUDED.counter`

	result, err := s.NamedExecContext(context.Background(), query, model)
	if err != nil {
		l.Err(err).Msg("p.db.NamedExecContext")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		l.Err(err).Msg("result.RowsAffected")
		return
	}

	if rows == 0 {
		l.Err(err).Msg("result.RowsAffected")
		return
	}
}

func (s *Store) Updates(ctx context.Context, req []memory.Metric) error {
	for _, v := range req {
		model := models.Metric{
			ID:      v.Name,
			Gauge:   v.Value,
			Counter: v.Delta,
		}

		query := `
		INSERT INTO tbl_metrics(id,gauge,counter)
		VALUES (:id, :gauge, :counter)
		ON CONFLICT (id) DO UPDATE
		SET
		    gauge = EXCLUDED.gauge,
		    counter = EXCLUDED.counter`

		result, err := s.NamedExecContext(context.Background(), query, model)
		if err != nil {
			return fmt.Errorf("p.DB.ExecContext: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("result.RowsAffected: %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("model: %+v; rows empty %w", model, errors.New("no rows affected"))
		}

	}

	return nil
}

func (s *Store) Ping() error {
	return s.storage.Ping()
}

func (s *Store) Close() error {
	return s.storage.Close()
}
