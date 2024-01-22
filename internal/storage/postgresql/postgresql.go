package postgresql

import (
	"context"
	"database/sql"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

type PsqlStorage struct {
	Conf *config.Config
	db   *sql.DB
	ctx  context.Context
}

func (store *PsqlStorage) Init(ctx context.Context) error {
	var err error
	store.ctx = ctx
	store.db, err = sql.Open("pgx", store.Conf.DbConnect)
	if err != nil {
		return err
	}
	_, err = store.db.QueryContext(store.ctx, "CREATE TABLE IF NOT EXISTS gauge (key varchar(100), value DOUBLE PRECISION);")
	if err != nil {
		return err
	}
	_, err = store.db.QueryContext(store.ctx, "CREATE TABLE IF NOT EXISTS counter (key varchar(100), value int8);")
	if err != nil {
		return err
	}
	return nil
}

func (store *PsqlStorage) Close() {
	store.db.Close()
}

func (store *PsqlStorage) AddGauge(key string, gauge float64) error {
	var err error
	var count int
	gaugeCountQuery := store.db.QueryRowContext(store.ctx, "SELECT COUNT(*) FROM gauge WHERE key=$1", key)
	if err = gaugeCountQuery.Err(); err != nil {
		logrus.Info(err)
		return err
	}
	err = gaugeCountQuery.Scan(&count)
	if err != nil {
		logrus.Info(err)
		return err
	}
	if count == 0 {
		_, err = store.db.ExecContext(store.ctx, "INSERT INTO gauge (key, value) VALUES ($1, $2)", key, gauge)
		if err != nil {
			logrus.Info(err)
			return err
		}
	} else {
		_, err = store.db.ExecContext(store.ctx, "UPDATE gauge SET value=$1 WHERE key=$2", gauge, key)
		if err != nil {
			logrus.Info(err)
			return err
		}
	}
	return nil
}

func (store *PsqlStorage) AddCounter(key string, counter int64) error {
	var err error
	var count int
	row := store.db.QueryRowContext(store.ctx, "SELECT COUNT(*) FROM counter WHERE key=$1", key)
	if err = row.Err(); err != nil {
		logrus.Info(err)
		return err
	}
	err = row.Scan(&count)
	if err != nil {
		logrus.Info(err)
		return err
	}
	if count == 0 {
		_, err = store.db.ExecContext(store.ctx, "INSERT INTO counter (key, value) VALUES ($1, $2)", key, counter)
		if err != nil {
			logrus.Info(err)
			return err
		}
	} else {
		var counterVal int64
		row2 := store.db.QueryRowContext(store.ctx, "SELECT value FROM counter WHERE key=$1", key)
		if err = row2.Err(); err != nil {
			logrus.Info(err)
			return err
		}
		err = row2.Scan(&counterVal)
		if err != nil {
			logrus.Info(err)
			return err
		}
		_, err = store.db.ExecContext(store.ctx, "UPDATE counter SET value=$1 WHERE key=$2", counterVal+counter, key)
		if err != nil {
			logrus.Info(err)
			return err
		}
	}

	return nil
}

func (store *PsqlStorage) GetCounter(key string) (int64, bool) {
	var v sql.NullInt64
	row := store.db.QueryRowContext(store.ctx, "SELECT value FROM counter WHERE key=$1", key)
	if err := row.Err(); err != nil {
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

func (store *PsqlStorage) GetGauge(key string) (float64, bool) {
	var v sql.NullFloat64
	row := store.db.QueryRowContext(store.ctx, "SELECT value FROM gauge WHERE key=$1", key)
	if err := row.Err(); err != nil {
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

func (store *PsqlStorage) GetAllCounters() map[string]int64 {
	counter := make(map[string]int64, 0)
	rows, err := store.db.QueryContext(store.ctx, "SELECT key, value FROM counter")
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
	return counter
}

func (store *PsqlStorage) GetAllGauge() map[string]float64 {
	counter := make(map[string]float64, 0)
	rows, err := store.db.QueryContext(store.ctx, "SELECT key, value FROM gauge")
	if err != nil {
		logrus.Info(err)
		return counter
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var val float64
		if err := rows.Scan(&key, &val); err != nil {
			logrus.Info(err)
			return counter
		}
		counter[key] = val
	}
	return counter
}

func (store *PsqlStorage) SaveToFile() error {
	return nil
}

func (store *PsqlStorage) Ping() error {
	return store.db.Ping()
}
