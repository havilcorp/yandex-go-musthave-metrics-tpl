package repositories

type IDataBase interface {
	Ping() error
}

type DataBase struct{ db IDataBase }

func NewDataBase(db IDataBase) *DataBase {
	return &DataBase{
		db: db,
	}
}

func (db *DataBase) Ping() error {
	return db.db.Ping()
}
