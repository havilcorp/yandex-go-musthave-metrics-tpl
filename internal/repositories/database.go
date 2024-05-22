// Package repositories репозиторий для проверки статуса подключения к бд
package repositories

type DataBaseSaver interface {
	Ping() error
}

type DataBase struct{ db DataBaseSaver }

func NewDataBase(db DataBaseSaver) *DataBase {
	return &DataBase{
		db: db,
	}
}

func (db *DataBase) Ping() error {
	return db.db.Ping()
}
