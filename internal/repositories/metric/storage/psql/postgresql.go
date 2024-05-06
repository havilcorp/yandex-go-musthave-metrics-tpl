// Package psql репозиторий для работы с метриками в базе данных
package psql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config/server"
	"github.com/sirupsen/logrus"
)

type PsqlStorage struct {
	db *sql.DB
}

// NewPsqlStorage инициализация хранилища в базе даннных + создание структуры
func NewPsqlStorage(conf *server.Config, db *sql.DB) (*PsqlStorage, error) {
	ctx := context.Background()
	psqlStorage := PsqlStorage{
		db: db,
	}
	if err := psqlStorage.Bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("init => %w", err)
	}
	return &psqlStorage, nil
}

// Bootstrap создание структуры
func (store *PsqlStorage) Bootstrap(ctx context.Context) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("bootstrap => %w", err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.Info(err)
		}
	}()
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS gauge (
			id SERIAL PRIMARY KEY,
			key varchar(100) UNIQUE NOT NULL, 
			value DOUBLE PRECISION NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS counter (
			id SERIAL PRIMARY KEY,
			key varchar(100) UNIQUE NOT NULL, 
			value bigint NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// AddGauge добавление метрики
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

// AddCounter добавление метрики
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

// AddGaugeBulk добавление метрики массивом в транзации
func (store *PsqlStorage) AddGaugeBulk(ctx context.Context, list []domain.Gauge) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("addGaugeBulk => %w", err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.Error(err)
		}
	}()
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

// AddCounterBulk добавление метрики массивом в транзации
func (store *PsqlStorage) AddCounterBulk(ctx context.Context, list []domain.Counter) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("addCounterBulk => %w", err)
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			logrus.Error(err)
		}
	}()
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

// GetGauge получение значения метрики
func (store *PsqlStorage) GetGauge(ctx context.Context, key string) (float64, error) {
	var v sql.NullFloat64
	var row *sql.Row
	var err error
	for _, sec := range []int{1, 3, 5} {
		row = store.db.QueryRowContext(ctx, "SELECT value FROM gauge WHERE key=$1", key)
		if err = row.Err(); err != nil {
			logrus.Error(err)
			time.Sleep(time.Duration(sec) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		logrus.Error(err)
		return 0, fmt.Errorf("getGauge => %w", err)
	}
	if err := row.Scan(&v); err != nil {
		logrus.Error(err)
		return 0, domain.ErrValueNotFound
	}
	if !v.Valid {
		return 0, domain.ErrValueNotFound
	}
	return v.Float64, nil
}

// GetCounter получение значения метрики
func (store *PsqlStorage) GetCounter(ctx context.Context, key string) (int64, error) {
	var v sql.NullInt64
	var row *sql.Row
	var err error
	for _, sec := range []int{1, 3, 5} {
		row = store.db.QueryRowContext(ctx, "SELECT value FROM counter WHERE key=$1", key)
		if err = row.Err(); err != nil {
			logrus.Error(err)
			time.Sleep(time.Duration(sec) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		logrus.Error(err)
		return 0, fmt.Errorf("getCounter => %w", err)
	}
	if err := row.Scan(&v); err != nil {
		logrus.Error(err)
		return 0, domain.ErrValueNotFound
	}
	if !v.Valid {
		return 0, domain.ErrValueNotFound
	}
	return v.Int64, nil
}

// GetAllGauge получение всех значений метрики
func (store *PsqlStorage) GetAllGauge(ctx context.Context) (map[string]float64, error) {
	gauge := make(map[string]float64, 0)
	var rows *sql.Rows
	var err error
	for _, sec := range []int{1, 3, 5} {
		rows, err = store.db.QueryContext(ctx, "SELECT key, value FROM gauge")
		if err != nil {
			logrus.Error(err)
			time.Sleep(time.Duration(sec) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		logrus.Error(err)
		return gauge, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	for rows.Next() {
		var key string
		var val float64
		if err := rows.Scan(&key, &val); err != nil {
			logrus.Error(err)
			return gauge, err
		}
		gauge[key] = val
	}
	if err := rows.Err(); err != nil {
		logrus.Error(err)
		return gauge, err
	}
	return gauge, nil
}

// GetAllCounters получение всех значений метрики
func (store *PsqlStorage) GetAllCounters(ctx context.Context) (map[string]int64, error) {
	counter := make(map[string]int64, 0)
	var rows *sql.Rows
	var err error
	for _, sec := range []int{1, 3, 5} {
		rows, err = store.db.QueryContext(ctx, "SELECT key, value FROM counter")
		if err != nil {
			logrus.Error(err)
			time.Sleep(time.Duration(sec) * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		logrus.Error(err)
		return counter, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Error(err)
		}
	}()
	for rows.Next() {
		var key string
		var val int64
		if err := rows.Scan(&key, &val); err != nil {
			logrus.Error(err)
			return counter, err
		}
		counter[key] = val
	}
	if err := rows.Err(); err != nil {
		logrus.Error(err)
		return counter, err
	}
	return counter, nil
}
