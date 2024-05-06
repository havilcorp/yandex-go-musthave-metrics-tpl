// Package metric фабрики хранилища
package metric

import (
	"database/sql/driver"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/internal/config/server"
	"github.com/sirupsen/logrus"
)

func TestMetricFactory(t *testing.T) {
	conf := server.NewServerConfig()
	var err error
	_, err = MetricFactory("memory", conf, nil)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = MetricFactory("file", conf, nil)
	if err != nil {
		t.Error(err)
		return
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logrus.Error(err)
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
	_, err = MetricFactory("psql", conf, db)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = MetricFactory("none", conf, nil)
	if err == nil {
		t.Error(err)
		return
	}
}
