package repositories

import (
	"testing"

	"github.com/havilcorp/yandex-go-musthave-metrics-tpl/mocks"
)

func TestDataBase_Ping(t *testing.T) {
	database := mocks.NewIDataBase(t)
	database.On("Ping").Return(nil)
	t.Run("Ping", func(t *testing.T) {
		db := NewDataBase(database)
		err := db.Ping()
		if err != nil {
			t.Errorf("Ping %v", err)
		}
	})
}