package services

import (
	"helpers/database"
	"log"
	proto "protos"
	"time"
)

type depositsStorageService struct {
	db      database.Database
	logger  *log.Logger
	queries map[string]string
}

type DepositsStorageService interface {
	Close()
	QueryGet(p interface{}, queryName string, params ...interface{}) error
	Exec(queryName string, params ...interface{}) error
	GetDepositsList(userId int64, page int64) ([]*proto.DepositItem, error)
	CreateDeposit(rsp *proto.DepositItem, userId int64, amount float32, currency string) error
}

// Deposits Storage Service
func NewDepositsStorageService(db database.Database, logger *log.Logger, queries map[string]string) DepositsStorageService {
	depositsStorage := &depositsStorageService{db, logger, queries}

	return depositsStorage
}

// Close the database connection
func (storage *depositsStorageService) Close() {
	storage.db.Close()
}

func (storage *depositsStorageService) QueryGet(p interface{}, queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.GetRecord(p, query, params...)
}

func (storage *depositsStorageService) Exec(queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.Exec(query, params...)
}

const itemsPerPage int64 = 20

func (storage *depositsStorageService) GetDepositsList(userId int64, page int64) ([]*proto.DepositItem, error) {
	list := []*proto.DepositItem{}
	query := storage.queries["GetDepositsList"]

	err := storage.db.GetRecords(&list, query, userId, itemsPerPage, page*itemsPerPage)
	return list, err
}

const depositActiveYears int = 2

func (storage *depositsStorageService) CreateDeposit(rsp *proto.DepositItem, userId int64, amount float32, currency string) error {
	timeNow := time.Now().UnixNano()
	validTill := time.Now().AddDate(depositActiveYears, 0, 0).UnixNano()
	return storage.QueryGet(rsp, "CreateDeposit", userId, amount, currency, timeNow, validTill)
}
