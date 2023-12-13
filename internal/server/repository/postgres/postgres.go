package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/1Asi1/metric-track.git/internal/server/models"
	"github.com/1Asi1/metric-track.git/internal/server/repository/memory"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
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

type Postgres struct {
	*sqlx.DB
	sync.RWMutex
	log zerolog.Logger
}

func New(cfg Config, log zerolog.Logger) (*Postgres, error) {
	l := cfg.Logger.With().Str("postgres", "New").Logger()

	db, err := sqlx.Connect("pgx", cfg.ConnURL)
	if err != nil {
		return &Postgres{}, fmt.Errorf("postgres connection error: %w", err)
	}
	l.Info().Msg("succeeded in connecting to postgres")

	db.SetConnMaxLifetime(cfg.MaxConnLifeTime)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)
	db.SetMaxOpenConns(cfg.MaxConn)

	if err = runMigrations(cfg.ConnURL); err != nil {
		return &Postgres{}, fmt.Errorf("runMigrations error: %w", err)
	}

	return &Postgres{DB: db, log: log}, nil
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

func (p *Postgres) Get(ctx context.Context) (map[string]memory.Type, error) {
	query := `
	SELECT
	    id,
	    gauge,
		counter
	FROM tbl_metrics
`
	var models []models.Metric
	if err := p.DB.SelectContext(ctx, &models, query); err != nil {
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

func (p *Postgres) GetOne(ctx context.Context, name string) (memory.Type, error) {
	query := `
	SELECT
	    gauge,
		counter
	FROM tbl_metrics
	WHERE id = $1
`
	var model memory.Type
	err := p.DB.GetContext(ctx, &model, query, name)
	if err != nil {
		return memory.Type{}, fmt.Errorf("GetOne: %w", err)
	}

	return model, nil
}

func (p *Postgres) Update(ctx context.Context, data map[string]memory.Type) {

	l := p.log.With().Str("postgres", "Update").Logger()

	name := fmt.Sprintf("%v", ctx.Value(Name("name")))
	model := models.Metric{
		ID:      name,
		Gauge:   data[name].Gauge,
		Counter: data[name].Counter,
	}
	p.Lock()
	chek, _ := p.GetOne(context.Background(), model.ID)
	if reflect.DeepEqual(chek, memory.Type{}) {
		query := `
		INSERT INTO tbl_metrics(id,gauge,counter)
		VALUES (:id, :gauge, :counter)`

		result, err := p.DB.NamedExecContext(context.Background(), query, model)
		if err != nil {
			l.Err(err).Msg("p.db.NamedExecContext")
		}

		rows, err := result.RowsAffected()
		if err != nil {
			l.Err(err).Msg("result.RowsAffected")
		}

		if rows == 0 {
			l.Err(err).Msg("result.RowsAffected")
		}
	} else {
		query := `
		UPDATE
		    tbl_metrics
		SET
		    gauge = :gauge,
		    counter = :counter
		WHERE
		    id = :id`

		result, err := p.DB.NamedExecContext(context.Background(), query, model)
		if err != nil {
			l.Err(err).Msg("p.db.NamedExecContext")
		}

		rows, err := result.RowsAffected()
		if err != nil {
			l.Err(err).Msg("result.RowsAffected")
		}

		if rows == 0 {
			l.Err(err).Msg("result.RowsAffected")
		}
	}
	p.Unlock()
}

func (p *Postgres) Ping() error {
	return p.DB.Ping()
}

func (p *Postgres) Updates(ctx context.Context, req []memory.Metric) error {
	l := p.log.With().Str("postgres", "Updates").Logger()

	for _, v := range req {
		_, err := p.GetOne(context.Background(), v.Name)
		if err == nil {
			if err = p.updateMetric(v); err != nil {
				l.Err(err).Msgf("p.createMetric metric value: %+v", v)
			}
		} else {
			if err = p.createMetric(v); err != nil {
				l.Err(err).Msgf("p.createMetric metric value: %+v", v)
			}
		}
	}

	return nil
}

func (p *Postgres) createMetric(req memory.Metric) error {
	model := models.Metric{
		ID:      req.Name,
		Gauge:   req.Value,
		Counter: req.Delta,
	}

	query := `
		INSERT INTO tbl_metrics(id,gauge,counter)
		VALUES (:id, :gauge, :counter)`
	p.Lock()
	result, err := p.DB.NamedExecContext(context.Background(), query, model)
	if err != nil {
		return fmt.Errorf("p.DB.ExecContext: %w", err)
	}
	p.Unlock()
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("model: %+v; rows empty %w", model, errors.New("no rows affected"))
	}
	return nil
}

func (p *Postgres) updateMetric(req memory.Metric) error {
	model := models.Metric{
		ID:      req.Name,
		Gauge:   req.Value,
		Counter: req.Delta,
	}

	query := `
		UPDATE
		    tbl_metrics
		SET
		    gauge = :gauge,
		    counter = :counter
		WHERE
		    id = :id`

	p.Lock()
	result, err := p.DB.NamedExecContext(context.Background(), query, model)
	if err != nil {
		return fmt.Errorf("p.DB.NamedExecContext: %w", err)
	}
	p.Unlock()

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("rows empty %w", errors.New("no rows affected"))
	}
	return nil
}
