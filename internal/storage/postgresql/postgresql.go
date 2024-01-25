package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"syscall"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

type PsqlStorage struct {
	Conf *config.Config
	db   *sql.DB
}

func (store *PsqlStorage) Init() error {
	ctx := context.Background()
	var err error
	store.db, err = sql.Open("pgx", store.Conf.DBConnect)
	if err != nil {
		return fmt.Errorf("init => %w", err)
	}
	for _, sec := range []int{1, 3, 5} {
		err = store.db.PingContext(ctx)
		if errors.Is(err, syscall.ECONNREFUSED) {
			time.Sleep(time.Duration(sec) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		store.db.Close()
		return fmt.Errorf("init => %w", err)
	}
	// if err, ok := err.(*pgconn.PgError); ok {
	// 	if err.Code == pgerrcode.ConnectionFailure {
	if err := store.Bootstrap(ctx); err != nil {
		return fmt.Errorf("init => %w", err)
	}
	return nil
}

func (store *PsqlStorage) Bootstrap(ctx context.Context) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("bootstrap => %w", err)
	}
	defer tx.Rollback()
	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gauge (
			id SERIAL PRIMARY KEY,
			key varchar(100) UNIQUE NOT NULL, 
			value DOUBLE PRECISION NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS counter (
			id SERIAL PRIMARY KEY,
			key varchar(100) UNIQUE NOT NULL, 
			value bigint NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return tx.Commit()
}

func (store *PsqlStorage) Close() {
	store.db.Close()
}

func (store *PsqlStorage) AddGauge(ctx context.Context, key string, gauge float64) error {
	_, err := store.db.ExecContext(ctx, `
		INSERT INTO gauge (key, value)
		VALUES($1, $2) 
		ON CONFLICT (key) 
		DO UPDATE SET value = $2;
	`, key, gauge)
	if err != nil {
		return fmt.Errorf("addGauge => %w", err)
	}
	return nil
}

func (store *PsqlStorage) AddCounter(ctx context.Context, key string, counter int64) error {
	_, err := store.db.ExecContext(ctx, `
		INSERT INTO counter (key, value)
		VALUES($1, $2) 
		ON CONFLICT (key) 
		DO UPDATE SET value = counter.value + $2;
	`, key, counter)
	if err != nil {
		return fmt.Errorf("addCounter => %w", err)
	}
	return nil
}

func (store *PsqlStorage) AddGaugeBulk(ctx context.Context, list []models.GaugeModel) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("addGaugeBulk => %w", err)
	}
	defer tx.Rollback()
	for _, model := range list {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
		`, model.Key, model.Value)
		if err != nil {
			return fmt.Errorf("addGaugeBulk => %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("addGaugeBulk => %w", err)
	}
	return nil
}

func (store *PsqlStorage) AddCounterBulk(ctx context.Context, list []models.CounterModel) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("addCounterBulk => %w", err)
	}
	defer tx.Rollback()
	for _, model := range list {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO counter (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = counter.value + $2;
		`, model.Key, model.Value)
		if err != nil {
			return fmt.Errorf("addCounterBulk => %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("addCounterBulk => %w", err)
	}
	return nil
}

func (store *PsqlStorage) GetCounter(ctx context.Context, key string) (int64, bool) {
	var v sql.NullInt64
	var row *sql.Row
	var err error
	for _, sec := range []int{1, 3, 5} {
		row = store.db.QueryRowContext(ctx, "SELECT value FROM counter WHERE key=$1", key)
		if err = row.Err(); err != nil {
			logrus.Info(err)
			time.Sleep(time.Duration(sec) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		logrus.Info(err)
		return 0, false
	}
	if err := row.Scan(&v); err != nil {
		logrus.Info(err)
		return 0, false
	}
	if !v.Valid {
		return 0, false
	}
	return v.Int64, true
}

func (store *PsqlStorage) GetGauge(ctx context.Context, key string) (float64, bool) {
	var v sql.NullFloat64
	var row *sql.Row
	var err error
	for _, sec := range []int{1, 3, 5} {
		row = store.db.QueryRowContext(ctx, "SELECT value FROM gauge WHERE key=$1", key)
		if err = row.Err(); err != nil {
			logrus.Info(err)
			time.Sleep(time.Duration(sec) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		logrus.Info(err)
		return 0, false
	}
	if err := row.Scan(&v); err != nil {
		logrus.Info(err)
		return 0, false
	}
	if !v.Valid {
		return 0, false
	}
	return v.Float64, true
}

func (store *PsqlStorage) GetAllCounters(ctx context.Context) map[string]int64 {
	counter := make(map[string]int64, 0)
	var rows *sql.Rows
	var err error
	for _, sec := range []int{1, 3, 5} {
		rows, err = store.db.QueryContext(ctx, "SELECT key, value FROM counter")
		if err != nil {
			logrus.Info(err)
			time.Sleep(time.Duration(sec) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		logrus.Info(err)
		return counter
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var val int64
		if err := rows.Scan(&key, &val); err != nil {
			logrus.Info(err)
			return counter
		}
		counter[key] = val
	}
	if err := rows.Err(); err != nil {
		logrus.Info(err)
		return counter
	}
	return counter
}

func (store *PsqlStorage) GetAllGauge(ctx context.Context) map[string]float64 {
	gauge := make(map[string]float64, 0)
	var rows *sql.Rows
	var err error
	for _, sec := range []int{1, 3, 5} {
		rows, err = store.db.QueryContext(ctx, "SELECT key, value FROM gauge")
		if err != nil {
			logrus.Info(err)
			time.Sleep(time.Duration(sec) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		logrus.Info(err)
		return gauge
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var val float64
		if err := rows.Scan(&key, &val); err != nil {
			logrus.Info(err)
			return gauge
		}
		gauge[key] = val
	}
	if err := rows.Err(); err != nil {
		logrus.Info(err)
		return gauge
	}
	return gauge
}

func (store *PsqlStorage) SaveToFile(ctx context.Context) error {
	return nil
}

func (store *PsqlStorage) Ping() error {
	return store.db.Ping()
}
