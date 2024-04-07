package psql

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/domain"
)

// func SetupTestDatabase() (testcontainers.Container, *sql.DB) {
// 	containerReq := testcontainers.ContainerRequest{
// 		Image:        "postgres:latest",
// 		ExposedPorts: []string{"5434/tcp"},
// 		WaitingFor:   wait.ForListeningPort("5434/tcp"),
// 		Env: map[string]string{
// 			"POSTGRES_DB":       "testdb",
// 			"POSTGRES_PASSWORD": "postgres",
// 			"POSTGRES_USER":     "postgres",
// 		},
// 	}
// 	dbContainer, _ := testcontainers.GenericContainer(
// 		context.Background(),
// 		testcontainers.GenericContainerRequest{
// 			ContainerRequest: containerReq,
// 			Started:          true,
// 		})

// 	host, _ := dbContainer.Host(context.Background())
// 	port, _ := dbContainer.MappedPort(context.Background(), "5434")
// 	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/postgres", host, port.Port())

// 	db, err := sql.Open("pgx", dbURI)
// 	if err != nil {
// 		logrus.Errorf("pgx init => %v", err)
// 		return nil, nil
// 	}
// 	defer db.Close()

// 	return dbContainer, db
// }

// func TestPsqlStorage_AddGauge(t *testing.T) {
// 	dbContainer, db := SetupTestDatabase()
// 	defer dbContainer.Terminate(context.Background())

// 	conf := config.NewConfig()
// 	dbStore, err := NewPsqlStorage(conf, db)
// 	if err != nil {
// 		return
// 	}
// }

func TestPsqlStorage_AddGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO gauge (key, value)
		VALUES($1, $2) 
		ON CONFLICT (key) 
		DO UPDATE SET value = $2;
	`)).WithArgs(
		"GAUGE",
		1.1,
	).WillReturnResult(driver.ResultNoRows)

	psqlStorage := PsqlStorage{
		db: db,
	}
	t.Run("AddGauge", func(t *testing.T) {
		err = psqlStorage.AddGauge(context.Background(), "GAUGE", 1.1)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestPsqlStorage_AddCounter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO counter (key, value)
		VALUES($1, $2) 
		ON CONFLICT (key) 
		DO UPDATE SET value = counter.value + $2;
	`)).WithArgs(
		"COUNTER",
		1,
	).WillReturnResult(driver.ResultNoRows)

	psqlStorage := PsqlStorage{
		db: db,
	}
	t.Run("AddCounter", func(t *testing.T) {
		err = psqlStorage.AddCounter(context.Background(), "COUNTER", 1)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestPsqlStorage_AddGaugeBulk(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

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
	t.Run("AddGaugeBulk", func(t *testing.T) {
		list := make([]domain.Gauge, 0)
		list = append(list, domain.Gauge{Key: "GAUGE1", Value: 1.1})
		list = append(list, domain.Gauge{Key: "GAUGE2", Value: 1.2})
		err = psqlStorage.AddGaugeBulk(context.Background(), list)
		if err != nil {
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
	t.Run("AddGaugeBulk", func(t *testing.T) {
		list := make([]domain.Counter, 0)
		list = append(list, domain.Counter{Key: "COUNTER1", Value: 1})
		list = append(list, domain.Counter{Key: "COUNTER2", Value: 2})
		err = psqlStorage.AddCounterBulk(context.Background(), list)
		if err != nil {
			t.Error(err)
		}
	})
}

func TestPsqlStorage_GetGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT value FROM gauge WHERE key=$1`)).
		WithArgs("GAUGE").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1.1))
	psqlStorage := PsqlStorage{
		db: db,
	}
	t.Run("GetGauge", func(t *testing.T) {
		val, err := psqlStorage.GetGauge(context.Background(), "GAUGE")
		if err != nil {
			t.Error(err)
		}
		if val != 1.1 {
			t.Error(errors.New("value not equil"))
		}
	})
}

func TestPsqlStorage_GetCounter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT value FROM counter WHERE key=$1`)).
		WithArgs("COUNTER").
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(1))
	psqlStorage := PsqlStorage{
		db: db,
	}
	t.Run("GetCounter", func(t *testing.T) {
		val, err := psqlStorage.GetCounter(context.Background(), "COUNTER")
		if err != nil {
			t.Error(err)
		}
		if val != 1 {
			t.Error(errors.New("value not equil"))
		}
	})
}

func TestPsqlStorage_GetAllGauge(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT key, value FROM gauge`)).
		WillReturnRows(sqlmock.NewRows([]string{"key", "value"}).
			AddRows([]driver.Value{"GAUGE1", 1.1}).
			AddRows([]driver.Value{"GAUGE2", 1.2}))
	psqlStorage := PsqlStorage{
		db: db,
	}
	t.Run("GetAllGauge", func(t *testing.T) {
		_, err := psqlStorage.GetAllGauge(context.Background())
		if err != nil {
			t.Error(err)
		}
	})
}

func TestPsqlStorage_GetAllCounters(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT key, value FROM counter`)).
		WillReturnRows(sqlmock.NewRows([]string{"key", "value"}).
			AddRows([]driver.Value{"COUNTER1", 1}).
			AddRows([]driver.Value{"COUNTER2", 2}))
	psqlStorage := PsqlStorage{
		db: db,
	}
	t.Run("GetAllCounters", func(t *testing.T) {
		_, err := psqlStorage.GetAllCounters(context.Background())
		if err != nil {
			t.Error(err)
		}
	})
}