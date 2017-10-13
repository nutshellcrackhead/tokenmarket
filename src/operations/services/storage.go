package services

import (
	"crypto/sha256"
	"encoding/base64"
	"helpers/database"
	"log"
	proto "protos"
	"strconv"
	"time"
)

type operationsStorageService struct {
	db      database.Database
	logger  *log.Logger
	queries map[string]string
	encoder *base64.Encoding
}

type OperationsStorageService interface {
	Close()
	CreatePayment(result *proto.CreatePaymentResponse, id int64, method string, amount float32, currencyName string, operationType string) error
	CreatePayout(id int64, method string, amount float32, currencyName string, account string, operationType string) error
	CreateAccountActivation(paymentId int64) error
	ValidatePayment(paymentId int64, token string, amount float32, currency string, method string) (int64, error)
	UpdatePayment(paymentId int64, status string) error
	UnlockAccount(paymentId int64) error
	TopUpWallet(operationId int64, status string) error
}

// Operations Storage Service
func NewOperationsStorageService(db database.Database, logger *log.Logger, queries map[string]string, encoder *base64.Encoding) OperationsStorageService {
	operationsStorage := &operationsStorageService{db, logger, queries, encoder}

	return operationsStorage
}

// Close the database connection
func (storage *operationsStorageService) Close() {
	storage.db.Close()
}

func (storage *operationsStorageService) CreatePayment(result *proto.CreatePaymentResponse, id int64, method string, amount float32, currencyName string, operationType string) error {
	timeNow := time.Now().UnixNano()
	token := storage.generatePaymentToken(timeNow)
	query := storage.queries["CreatePayment"]
	return storage.db.GetRecord(result, query, id, method, timeNow, amount, currencyName, token, operationType)
}

func (storage *operationsStorageService) CreatePayout(id int64, method string, amount float32, currencyName string, account string, operationType string) error {
	timeNow := time.Now().UnixNano()
	token := storage.generatePaymentToken(timeNow)
	query := storage.queries["CreatePayout"]
	return storage.db.Exec(query, id, method, timeNow, amount, currencyName, token, account, operationType)
}

func (storage *operationsStorageService) CreateAccountActivation(paymentId int64) error {
	query := storage.queries["CreateAccountActivation"]
	return storage.db.Exec(query, paymentId)
}

func (storage *operationsStorageService) ValidatePayment(paymentId int64, token string, amount float32, currency string, method string) (int64, error) {
	result := int64(0)
	query := storage.queries["ValidatePayment"]
	err := storage.db.GetRecord(&result, query, paymentId, token, amount, currency, method)
	return result, err
}

func (storage *operationsStorageService) UpdatePayment(paymentId int64, status string) error {
	query := storage.queries["UpdatePayment"]
	return storage.db.Exec(query, paymentId, status)
}

func (storage *operationsStorageService) UnlockAccount(paymentId int64) error {
	query := storage.queries["UnlockAccount"]
	return storage.db.Exec(query, paymentId)
}

func (storage *operationsStorageService) generatePaymentToken(timeNow int64) string {
	// generate token
	stringifiedTimeNow := strconv.FormatInt(timeNow, 10)

	sessionToken := sha256.New()
	sessionToken.Write([]byte(stringifiedTimeNow))
	shaToken := sessionToken.Sum(nil)

	return storage.encoder.EncodeToString(shaToken[:])
}

func (storage *operationsStorageService) TopUpWallet(operationId int64, status string) error {
	if status != "success" {
		return nil
	}

	query := storage.queries["TopUpWallet"]

	return storage.db.Exec(query, operationId)
}
