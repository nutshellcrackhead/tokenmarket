package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	microContext "golang.org/x/net/context"
	"net/http"
	protoSchema "protos"
)

type AuthServerController interface {
	LoginV1Auth(context echo.Context) error
	LogoutV1Auth(context echo.Context) error
	RegisterV1Auth(context echo.Context) error
}

type authController struct {
	microservice protoSchema.AuthClient
}

func NewAuthController(microservice protoSchema.AuthClient) AuthServerController {
	return &authController{microservice: microservice}
}

func (controller *authController) LoginV1Auth(context echo.Context) (err error) {
	loginRequest := &protoSchema.AuthLoginRequest{}

	if err = context.Bind(loginRequest); err != nil {
		fmt.Println(err)
		return context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "wrong_username_or_password"})
	}

	authResponse, err := controller.microservice.Login(microContext.TODO(), loginRequest)

	if err != nil {
		fmt.Println(err)
		return context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "wrong_username_or_password"})
	}

	return context.JSON(http.StatusOK, authResponse)
}

func (controller *authController) LogoutV1Auth(context echo.Context) (err error) {
	id := context.Get("userId").(int64)
	token := context.Get("token").(string)
	logoutRequest := &protoSchema.AuthLogoutRequest{Id: &id, Token: &token}

	_, err = controller.microservice.Logout(microContext.TODO(), logoutRequest)

	if err != nil {
		return context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "wrong_data"})
	}

	return context.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (controller *authController) RegisterV1Auth(context echo.Context) (err error) {
	registerRequest := &protoSchema.AuthRegisterRequest{}

	if err = context.Bind(registerRequest); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "wrong_username_or_email"})
		return err
	}

	registerResponse, err := controller.microservice.Register(microContext.TODO(), registerRequest)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "wrong_username_or_email"})
		return err
	}

	context.JSON(http.StatusOK, &registerResponse)
	return err
}
