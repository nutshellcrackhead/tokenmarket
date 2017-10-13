package services

import (
	"helpers/database"
	"log"
	proto "protos"
)

type referralStorageService struct {
	db      database.Database
	logger  *log.Logger
	queries map[string]string
}

type ReferralStorageService interface {
	Close()
	GetReferralStatus(p *[]*proto.ReferralPartner, userId int64, page int64) error
}

// Referral Storage Service
func NewReferralStorageService(db database.Database, logger *log.Logger, queries map[string]string) ReferralStorageService {
	referralStorage := &referralStorageService{db, logger, queries}

	return referralStorage
}

// Close the database connection
func (storage *referralStorageService) Close() {
	storage.db.Close()
}

const itemsPerPage int64 = 100

func (storage *referralStorageService) GetReferralStatus(p *[]*proto.ReferralPartner, userId int64, page int64) error {
	query := storage.queries["GetReferrers"]

	return storage.db.GetRecords(p, query, userId, itemsPerPage, page*itemsPerPage)
}
