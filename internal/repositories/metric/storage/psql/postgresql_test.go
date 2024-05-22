package psql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config/server"
	"github.com/sirupsen/logrus"
)

func TestPsqlStorage_AddGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Info(err)
		}
	}()

	sql := `
		INSERT INTO gauge (key, value)
		VALUES($1, $2) 
		ON CONFLICT (key) 
		DO UPDATE SET value = $2;
	`

	mock.ExpectExec(regexp.QuoteMeta(sql)).WithArgs("GAUGE", 1.1).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(regexp.QuoteMeta(sql)).WithArgs("GAUGE", 2.1).WillReturnError(errors.New("Error db"))

	psqlStorage := PsqlStorage{
		db: db,
	}
	err = psqlStorage.AddGauge(context.Background(), "GAUGE", 1.1)
	if err != nil {
		t.Error(err)
	}
	err = psqlStorage.AddGauge(context.Background(), "GAUGE", 2.1)
	if err.Error() != "Error db" {
		t.Error(err)
	}
}

func TestPsqlStorage_AddCounter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Info(err)
		}
	}()

	sql := `
		INSERT INTO counter (key, value)
		VALUES($1, $2) 
		ON CONFLICT (key) 
		DO UPDATE SET value = counter.value + $2;
	`

	mock.ExpectExec(regexp.QuoteMeta(sql)).WithArgs("COUNTER", 1).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(regexp.QuoteMeta(sql)).WithArgs("COUNTER", 2).WillReturnError(errors.New("Error db"))

	psqlStorage := PsqlStorage{
		db: db,
	}
	err = psqlStorage.AddCounter(context.Background(), "COUNTER", 1)
	if err != nil {
		t.Error(err)
	}
	err = psqlStorage.AddCounter(context.Background(), "COUNTER", 2)
	if err.Error() != "Error db" {
		t.Error(err)
	}
}

func TestPsqlStorage_AddGaugeBulk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("1", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
		`)).WithArgs(
			"GAUGE1",
			1.1,
		).WillReturnResult(driver.ResultNoRows)

		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
		`)).WithArgs(
			"GAUGE2",
			1.2,
		).WillReturnResult(driver.ResultNoRows)

		mock.ExpectCommit()

		psqlStorage := PsqlStorage{
			db: db,
		}
		list := make([]domain.Gauge, 0)
		list = append(list, domain.Gauge{Key: "GAUGE1", Value: 1.1})
		list = append(list, domain.Gauge{Key: "GAUGE2", Value: 1.2})
		err = psqlStorage.AddGaugeBulk(context.Background(), list)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("2", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.New("ExpectBegin"))
		psqlStorage := PsqlStorage{db: db}
		list := make([]domain.Gauge, 0)
		err = psqlStorage.AddGaugeBulk(context.Background(), list)
		if err.Error() != "ExpectBegin" {
			t.Error(err)
		}
	})
	t.Run("3", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
		`)).WithArgs(
			"GAUGE",
			2.1,
		).WillReturnError(errors.New("ExpectExec"))

		psqlStorage := PsqlStorage{db: db}
		list := make([]domain.Gauge, 0)
		list = append(list, domain.Gauge{Key: "GAUGE", Value: 2.1})
		err = psqlStorage.AddGaugeBulk(context.Background(), list)
		if err.Error() != "ExpectExec" {
			t.Error(err)
		}
	})
	t.Run("4", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO gauge (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = $2;
		`)).WithArgs("GAUGE", 1.1).WillReturnResult(driver.ResultNoRows)
		psqlStorage := PsqlStorage{db: db}
		list := make([]domain.Gauge, 0)
		list = append(list, domain.Gauge{Key: "GAUGE", Value: 1.1})
		mock.ExpectCommit().WillReturnError(errors.New("ExpectCommit"))
		err = psqlStorage.AddGaugeBulk(context.Background(), list)
		if !errors.Is(err, sql.ErrTxDone) {
			t.Error(err)
		}
	})
}

func TestPsqlStorage_AddCounterBulk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	t.Run("1", func(t *testing.T) {
		mock.ExpectBegin()

		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO counter (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = counter.value + $2;
		`)).WithArgs(
			"COUNTER1",
			1,
		).WillReturnResult(driver.ResultNoRows)

		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO counter (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = counter.value + $2;
		`)).WithArgs(
			"COUNTER2",
			2,
		).WillReturnResult(driver.ResultNoRows)

		mock.ExpectCommit()

		psqlStorage := PsqlStorage{
			db: db,
		}
		list := make([]domain.Counter, 0)
		list = append(list, domain.Counter{Key: "COUNTER1", Value: 1})
		list = append(list, domain.Counter{Key: "COUNTER2", Value: 2})
		err = psqlStorage.AddCounterBulk(context.Background(), list)
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("2", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.New("ExpectBegin"))
		psqlStorage := PsqlStorage{db: db}
		list := make([]domain.Counter, 0)
		err = psqlStorage.AddCounterBulk(context.Background(), list)
		if err.Error() != "ExpectBegin" {
			t.Error(err)
		}
	})
	t.Run("3", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO counter (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = counter.value + $2;
		`)).WithArgs(
			"COUNTER",
			2,
		).WillReturnError(errors.New("ExpectExec"))
		psqlStorage := PsqlStorage{db: db}
		list := make([]domain.Counter, 0)
		list = append(list, domain.Counter{Key: "COUNTER", Value: 2})
		err = psqlStorage.AddCounterBulk(context.Background(), list)
		if err.Error() != "ExpectExec" {
			t.Error(err)
		}
	})
	t.Run("4", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`
			INSERT INTO counter (key, value)
			VALUES($1, $2) 
			ON CONFLICT (key) 
			DO UPDATE SET value = counter.value + $2;
		`)).WithArgs("COUNTER", 1).WillReturnResult(driver.ResultNoRows)
		psqlStorage := PsqlStorage{db: db}
		list := make([]domain.Counter, 0)
		list = append(list, domain.Counter{Key: "COUNTER", Value: 1})
		mock.ExpectCommit().WillReturnError(errors.New(""))
		err = psqlStorage.AddCounterBulk(context.Background(), list)
		if !errors.Is(err, sql.ErrTxDone) {
			t.Error(err)
		}
	})
}

func TestPsqlStorage_GetGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Info(err)
		}
	}()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT value FROM gauge WHERE key=$1`)).
		WithArgs("GAUGE").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1.1))
	psqlStorage := PsqlStorage{
		db: db,
	}
	val, err := psqlStorage.GetGauge(context.Background(), "GAUGE")
	if err != nil {
		t.Error(err)
	}
	if val != 1.1 {
		t.Error(errors.New("value not equil"))
	}
}

func TestPsqlStorage_GetCounter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Info(err)
		}
	}()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT value FROM counter WHERE key=$1`)).
		WithArgs("COUNTER").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1))
	psqlStorage := PsqlStorage{
		db: db,
	}
	val, err := psqlStorage.GetCounter(context.Background(), "COUNTER")
	if err != nil {
		t.Error(err)
	}
	if val != 1 {
		t.Error(errors.New("value not equil"))
	}
}

func TestPsqlStorage_GetAllGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Info(err)
		}
	}()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT key, value FROM gauge`)).
		WillReturnRows(sqlmock.NewRows([]string{"key", "value"}).
			AddRows([]driver.Value{"GAUGE1", 1.1}).
			AddRows([]driver.Value{"GAUGE2", 1.2}))
	psqlStorage := PsqlStorage{
		db: db,
	}
	_, err = psqlStorage.GetAllGauge(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestPsqlStorage_GetAllCounters(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Info(err)
		}
	}()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT key, value FROM counter`)).
		WillReturnRows(sqlmock.NewRows([]string{"key", "value"}).
			AddRows([]driver.Value{"COUNTER1", 1}).
			AddRows([]driver.Value{"COUNTER2", 2}))
	psqlStorage := PsqlStorage{
		db: db,
	}
	_, err = psqlStorage.GetAllCounters(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestPsqlStorage_Bootstrap(t *testing.T) {
	conf := server.NewServerConfig()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Info(err)
		}
	}()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`
		CREATE TABLE IF NOT EXISTS gauge (
			id SERIAL PRIMARY KEY,
			key varchar(100) UNIQUE NOT NULL, 
			value DOUBLE PRECISION NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)).WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec(regexp.QuoteMeta(`
		CREATE TABLE IF NOT EXISTS counter (
			id SERIAL PRIMARY KEY,
			key varchar(100) UNIQUE NOT NULL, 
			value bigint NOT NULL,
			created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)).WillReturnResult(driver.ResultNoRows)
	mock.ExpectCommit()
	_, err = NewPsqlStorage(conf, db)
	if err != nil {
		t.Error(err)
	}
}
