package main

import (
	"github.com/labstack/echo"
	micro "github.com/micro/go-micro"
	protoSchema "protos"

	"api/controllers"
	"api/middleware"
	"api/server"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/broker/kafka"
	"github.com/micro/go-plugins/registry/zookeeper"
	"github.com/micro/go-plugins/transport/tcp"
)

func initMicroservice(name string) *micro.Service {
	registryAddress := registry.Addrs("token-kafka")
	brokerAddress := broker.Addrs("token-kafka:9092")

	registry := zookeeper.NewRegistry(registryAddress)
	transport := tcp.NewTransport()
	broker := kafka.NewBroker(brokerAddress)

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name(name),
		micro.Registry(registry),
		micro.Transport(transport),
		micro.Broker(broker),
	)

	return &service
}

func initProfileMicroservice() protoSchema.ProfileClient {
	service := *initMicroservice("profile.client")

	// Create new greeter client
	profile := protoSchema.NewProfileClient("profile", service.Client())

	return profile
}

func initAuthMicroservice() protoSchema.AuthClient {
	service := *initMicroservice("auth.client")

	// Create new greeter client
	auth := protoSchema.NewAuthClient("auth", service.Client())

	return auth
}

func initWalletMicroservice() protoSchema.WalletClient {
	service := *initMicroservice("wallet.client")

	// Create new greeter client
	wallet := protoSchema.NewWalletClient("wallet", service.Client())

	return wallet
}

func initReferralMicroservice() protoSchema.ReferralClient {
	service := *initMicroservice("referral.client")

	// Create new greeter client
	referral := protoSchema.NewReferralClient("referral", service.Client())

	return referral
}

func initOperationsMicroservice() protoSchema.OperationsClient {
	service := *initMicroservice("operations.client")

	// Create new greeter client
	operations := protoSchema.NewOperationsClient("operations", service.Client())

	return operations
}

func initDepositsMicroservice() protoSchema.DepositsClient {
	service := *initMicroservice("deposits.client")

	// Create new greeter client
	deposits := protoSchema.NewDepositsClient("deposits", service.Client())

	return deposits
}

func main() {
	authMicroservice := initAuthMicroservice()
	walletMicroservice := initWalletMicroservice()
	profileMicroservice := initProfileMicroservice()
	depositsMicroservice := initDepositsMicroservice()
	referralMicroservice := initReferralMicroservice()
	operationsMicroservice := initOperationsMicroservice()

	authControllers := controllers.NewAuthController(authMicroservice)
	walletControllers := controllers.NewWalletController(walletMicroservice)
	profileControllers := controllers.NewProfileController(profileMicroservice)
	depositsControllers := controllers.NewDepositsController(depositsMicroservice)
	referralControllers := controllers.NewReferralController(referralMicroservice)
	operationsControllers := controllers.NewOperationsController(operationsMicroservice)

	echoInstance := echo.New()
	authMiddleware := middleware.NewAuth(authMicroservice)

	serverInstance := server.NewServer(echoInstance, profileControllers,
		authControllers, walletControllers, referralControllers, operationsControllers,
		depositsControllers,
		authMiddleware)

	serverInstance.InitEndpoints()
	serverInstance.Run()
}
