package main

import (
	"fmt"
	micro "github.com/micro/go-micro"
	protoSchema "protos"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/broker"
	"auth/handlers"
	"auth/services"
	"github.com/micro/go-plugins/broker/kafka"
	"github.com/micro/go-plugins/registry/zookeeper"
	"github.com/micro/go-plugins/transport/tcp"
	"helpers/database"
	"log"
	"time"
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
		micro.Name("auth"),
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

const sessionTimeToLive time.Duration = 10

// Inits the shared profile storage (connection to database)
func initAuthStorage() services.AuthStorageService {
	db := database.Open()
	cache := database.NewCache()
	logger := &log.Logger{}
	queries := services.GetQueries()
	storage := services.NewAuthStorageService(db, cache, logger, queries, sessionTimeToLive)
	return storage
}

func main() {
	// Storage service init
	storage := initAuthStorage()
	defer storage.Close()

	service := initService()

	fmt.Println("Auth microservice is up and running.")

	// Register handler
	protoSchema.RegisterAuthHandler(service.Server(), handlers.NewAuthHandler(storage))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
