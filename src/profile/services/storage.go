package services

import (
	"helpers/database"
	"log"
	proto "protos"
)

type profileStorageService struct {
	db      database.Database
	logger  *log.Logger
	queries map[string]string
}

type ProfileStorageService interface {
	Close()
	QueryGet(p interface{}, queryName string, params ...interface{}) error
	QuerySelect(p *[]*proto.City, queryName string, params ...interface{}) error
}

// Profile Storage Service
func NewProfileStorageService(db database.Database, logger *log.Logger, queries map[string]string) ProfileStorageService {
	profileStorage := &profileStorageService{db, logger, queries}

	return profileStorage
}

// Close the database connection
func (storage *profileStorageService) Close() {
	storage.db.Close()
}

func (storage *profileStorageService) QueryGet(p interface{}, queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.GetRecord(p, query, params...)
}

func (storage *profileStorageService) QuerySelect(p *[]*proto.City, queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.GetRecords(p, query, params...)
}
