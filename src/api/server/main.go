package server

import (
	"api/controllers"
	"api/middleware"
	"github.com/labstack/echo"
	"path"
	"runtime"
)

type Server interface {
	Run()
	InitEndpoints()
}

type apiServer struct {
	echo                  *echo.Echo
	profileControllers    controllers.ProfileServerController
	authControllers       controllers.AuthServerController
	walletControllers     controllers.WalletServerController
	referralControllers   controllers.ReferralServerController
	operationsControllers controllers.OperationsServerController
	depositsControllers   controllers.DepositsServerController
	auth                  middleware.AuthMiddleware
}

func NewServer(echoInstance *echo.Echo,
	profileControllers controllers.ProfileServerController,
	authControllers controllers.AuthServerController,
	walletControllers controllers.WalletServerController,
	referralControllers controllers.ReferralServerController,
	operationsControllers controllers.OperationsServerController,
	depositsControllers controllers.DepositsServerController,
	auth middleware.AuthMiddleware) Server {
	return &apiServer{
		echo:                  echoInstance,
		profileControllers:    profileControllers,
		authControllers:       authControllers,
		walletControllers:     walletControllers,
		referralControllers:   referralControllers,
		operationsControllers: operationsControllers,
		depositsControllers:   depositsControllers,
		auth:                  auth,
	}
}

func (s *apiServer) InitEndpoints() {
	s.InitAuthServerEndpoints()
	s.InitProfileServerEndpoints()
	s.InitWalletServerEndpoints()
	s.InitReferralServerEndpoints()
	s.InitOperationsServerEndpoints()
	s.InitDepositsServerEndpoints()
}

func (s *apiServer) Run() {
	s.echo.Logger.Fatal(s.echo.StartTLS(":8080", getCertificatePath("cert.pem"), getCertificatePath("key.pem")))
}

func getCertificatePath(certName string) string {
	_, filename, _, _ := runtime.Caller(1)
	pwd := path.Join(path.Dir(filename))
	return pwd + "/" + certName
}

func (s *apiServer) InitProfileServerEndpoints() {
	profileGroup := s.echo.Group("/v1/profile", s.auth.Auth)

	profileGroup.GET("/my", s.profileControllers.GetV1ProfileMy)
	profileGroup.PATCH("/my", s.profileControllers.PatchV1ProfileMy)
	profileGroup.GET("/cities", s.profileControllers.CitiesV1Profile)
}

func (s *apiServer) InitAuthServerEndpoints() {
	authGroup := s.echo.Group("/v1/auth")

	authGroup.POST("/register", s.authControllers.RegisterV1Auth)
	authGroup.POST("/login", s.authControllers.LoginV1Auth)
	authGroup.GET("/logout", s.authControllers.LogoutV1Auth, s.auth.Auth)
}

func (s *apiServer) InitWalletServerEndpoints() {
	walletGroup := s.echo.Group("/v1/wallet", s.auth.Auth)

	walletGroup.GET("/my", s.walletControllers.GetV1WalletMy)
	walletGroup.GET("/my/operations", s.walletControllers.OperationsV1WalletMy)
}

func (s *apiServer) InitReferralServerEndpoints() {
	referralGroup := s.echo.Group("/v1/referral", s.auth.Auth)

	referralGroup.GET("/my", s.referralControllers.GetV1ReferralMy)
}

func (s *apiServer) InitOperationsServerEndpoints() {
	operationsGroup := s.echo.Group("/v1/operations")

	paymentsGroup := operationsGroup.Group("/payments")
	paymentsGroup.POST("/new", s.operationsControllers.NewV1Operations, s.auth.Auth)
	paymentsGroup.POST("/activate", s.operationsControllers.ActivateV1Operations, s.auth.Auth)
	paymentsGroup.POST("/status", s.operationsControllers.StatusV1Operations)

	operationsGroup.POST("/payout/new", s.operationsControllers.NewV1PayoutOperations, s.auth.Auth)
}

func (s *apiServer) InitDepositsServerEndpoints() {
	depositsGroup := s.echo.Group("/v1/deposits")

	depositsGroup.GET("/my", s.depositsControllers.GetV1DepositsMy, s.auth.Auth)
	depositsGroup.POST("/new", s.depositsControllers.NewV1Deposits, s.auth.Auth)
}
