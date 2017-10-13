package main

import (
	"fmt"
	micro "github.com/micro/go-micro"
	protoSchema "protos"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/broker/kafka"
	"github.com/micro/go-plugins/registry/zookeeper"
	"github.com/micro/go-plugins/transport/tcp"
	"helpers/database"
	"log"
	"profile/handlers"
	"profile/models"
	"profile/services"
)

// inits microservice
func initService() micro.Service {
	registryAddress := registry.Addrs("token-kafka")
        brokerAddress := broker.Addrs("token-kafka:9092")

	registry := zookeeper.NewRegistry(registryAddress)
	transport := tcp.NewTransport()
        broker := kafka.NewBroker(brokerAddress)

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("profile"),
		micro.Version("latest"),
		micro.Registry(registry),
		micro.Transport(transport),
		micro.Broker(broker),
	)

	// Init will parse the command line flags. Any flags set will
	// override the above settings. Options defined here will
	// override anything set on the command line.
	service.Init()

	return service
}

// Inits the shared profile storage (connection to database)
func initProfileStorage() services.ProfileStorageService {
	db := database.Open()
	logger := &log.Logger{}
	queries := services.GetQueries()
	storage := services.NewProfileStorageService(db, logger, queries)
	return storage
}

func main() {
	// Storage service init
	storage := initProfileStorage()
	defer storage.Close()

	service := initService()

	var GetProfileInstance = func(id int64) models.Profile {
		profileData := models.NewProfileData(storage)
		profileInstance := models.NewProfile(profileData, id)

		return profileInstance
	}

	fmt.Println("Profile microservice is up and running.")

	// Register handler
	protoSchema.RegisterProfileHandler(service.Server(), handlers.NewProfileHandler(GetProfileInstance))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
