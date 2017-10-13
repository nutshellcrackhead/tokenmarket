package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	microContext "golang.org/x/net/context"
	"net/http"
	protoSchema "protos"
	"strconv"
)

type WalletServerController interface {
	GetV1WalletMy(context echo.Context) error
	OperationsV1WalletMy(context echo.Context) error
}

type walletController struct {
	microservice protoSchema.WalletClient
}

func NewWalletController(microservice protoSchema.WalletClient) WalletServerController {
	return &walletController{microservice: microservice}
}

func (controller *walletController) GetV1WalletMy(context echo.Context) (err error) {
	id := context.Get("userId").(int64)
	walletStatus, err := controller.microservice.Status(microContext.TODO(), &protoSchema.WalletStatusRequest{Id: &id})

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, walletStatus)
	return
}

func (controller *walletController) OperationsV1WalletMy(context echo.Context) (err error) {
	var page int64

	id := context.Get("userId").(int64)
	pageString := context.QueryParam("page")
	parsedPage, parsePageErr := strconv.ParseInt(pageString, 10, 64)

	if parsePageErr == nil {
		page = parsedPage
	}

	operations, err := controller.microservice.OperationsList(microContext.TODO(),
		&protoSchema.WalletOperationsListRequest{Id: &id, Page: &page})

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, operations)
	return
}
