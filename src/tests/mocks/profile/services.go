package profile

import (
	"github.com/stretchr/testify/mock"
)

type ProfileStorageServiceMock struct {
	mock.Mock
}

func (storage *ProfileStorageServiceMock) Close() {
	storage.Called()
}

func (storage *ProfileStorageServiceMock) QueryGet(p interface{}, queryName string, params ...interface{}) {
	arguments := []interface{}{
		p,
		queryName,
	}

	storage.Called(append(arguments, params...)...)
}
