package services

import (
	"helpers/database"
	"log"
	"time"
)

type bonusStorageService struct {
	db      database.Database
	logger  *log.Logger
	queries map[string]string
}

type BonusStorageService interface {
	Close()
	QueryGet(p interface{}, queryName string, params ...interface{}) error
	Exec(queryName string, params ...interface{}) error
	PayDepositBonuses()
	PayReferralBonuses()
}

// Bonus Storage Service
func NewBonusStorageService(db database.Database, logger *log.Logger, queries map[string]string) BonusStorageService {
	bonusStorage := &bonusStorageService{db, logger, queries}

	return bonusStorage
}

// Close the database connection
func (storage *bonusStorageService) Close() {
	storage.db.Close()
}

func (storage *bonusStorageService) QueryGet(p interface{}, queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.GetRecord(p, query, params...)
}

func (storage *bonusStorageService) Exec(queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.Exec(query, params...)
}

func (storage *bonusStorageService) PayDepositBonuses() {
	storage.logger.Println("It's time to pay deposits.")

	timeNow := time.Now().UnixNano()

	if err := storage.Exec("PutDepositsBonuses", timeNow); err != nil {
		storage.logger.Println(err)
		return
	}

	storage.logger.Println("Success with deposits.")
}

func (storage *bonusStorageService) PayReferralBonuses() {
	storage.logger.Println("It's time to pay referral bonuses.")

	timeNow := time.Now()

	timeNowUnixNano := timeNow.UnixNano()
	weekAgoUnixNano := timeNow.AddDate(0, 0, -7).UnixNano()

	if err := storage.Exec("PutReferralBonuses", timeNowUnixNano, weekAgoUnixNano); err != nil {
		storage.logger.Println(err)
		return
	}

	storage.logger.Println("Success with referrals.")
}
