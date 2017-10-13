package services

import (
	"helpers/database"
	"log"
	proto "protos"
)

type walletStorageService struct {
	db      database.Database
	logger  *log.Logger
	queries map[string]string
}

type WalletStorageService interface {
	Close()
	QueryGet(p interface{}, queryName string, params ...interface{}) error
	Exec(queryName string, params ...interface{}) error
	GetWalletStatus(id int64) ([]*proto.WalletStatus, error)
	GetOperationsList(id int64, page int64) ([]*proto.WalletOperation, error)
}

// Wallet Storage Service
func NewWalletStorageService(db database.Database, logger *log.Logger, queries map[string]string) WalletStorageService {
	walletStorage := &walletStorageService{db, logger, queries}

	return walletStorage
}

// Close the database connection
func (storage *walletStorageService) Close() {
	storage.db.Close()
}

func (storage *walletStorageService) QueryGet(p interface{}, queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.GetRecord(p, query, params...)
}

func (storage *walletStorageService) Exec(queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.Exec(query, params...)
}

func (storage *walletStorageService) GetWalletStatus(id int64) ([]*proto.WalletStatus, error) {
	status := []*proto.WalletStatus{}
	query := storage.queries["GetWalletStatus"]
	err := storage.db.GetRecords(&status, query, id)

	return status, err
}

const operationsPerPage int64 = 20

func (storage *walletStorageService) GetOperationsList(id int64, page int64) ([]*proto.WalletOperation, error) {
	status := []*proto.WalletOperation{}
	query := storage.queries["GetUserOperations"]
	err := storage.db.GetRecords(&status, query, id, operationsPerPage, page*operationsPerPage)

	return status, err
}
