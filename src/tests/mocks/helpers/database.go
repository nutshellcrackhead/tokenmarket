package helpers

import (
	"github.com/stretchr/testify/mock"
)

type DatabaseMock struct {
	mock.Mock
}

func (db *DatabaseMock) Close() {
	db.Called()
}

func (db *DatabaseMock) GetRecord(instance interface{}, query string, params ...interface{}) {
	arguments := []interface{}{
		instance,
		query,
	}

	db.Called(append(arguments, params...)...)
}
