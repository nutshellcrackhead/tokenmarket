package main

import (
	"bonus/services"
	"fmt"
	"github.com/robfig/cron"
	"helpers/database"
	"log"
	"os"
	"os/signal"
	"time"
)

// Inits the shared profile storage (connection to database)
func initBonusStorage() services.BonusStorageService {
	db := database.Open()
	logger := log.New(os.Stdout, "", log.LstdFlags)
	queries := services.GetQueries()
	storage := services.NewBonusStorageService(db, logger, queries)
	return storage
}

func main() {
	// Storage service init
	storage := initBonusStorage()

	location, _ := time.LoadLocation("UTC")

	c := cron.NewWithLocation(location)

	// 4am UTC every Monday deposit bonuses
	c.AddFunc("0 0 4 * * 1", func() {
		storage.PayDepositBonuses()
	})

	// 4:15am UTC every Monday referral bonuses
	c.AddFunc("0 15 4 * * 1", func() {
		storage.PayReferralBonuses()
	})

	c.Start()

	fmt.Println("Bonus microservice is up and running.")

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
