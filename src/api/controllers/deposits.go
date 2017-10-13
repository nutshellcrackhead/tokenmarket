package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	microContext "golang.org/x/net/context"
	"net/http"
	protoSchema "protos"
	"strconv"
)

type DepositsServerController interface {
	GetV1DepositsMy(context echo.Context) error
	NewV1Deposits(context echo.Context) error
}

type depositsController struct {
	microservice protoSchema.DepositsClient
}

func NewDepositsController(microservice protoSchema.DepositsClient) DepositsServerController {
	return &depositsController{microservice: microservice}
}

func (controller *depositsController) GetV1DepositsMy(context echo.Context) (err error) {
	page := int64(0)

	id := context.Get("userId").(int64)
	pageString := context.QueryParam("page")
	parsedPage, parsePageErr := strconv.ParseInt(pageString, 10, 64)

	if parsePageErr == nil {
		page = parsedPage
	}

	depositsStatus, err := controller.microservice.GetDeposits(microContext.TODO(),
		&protoSchema.GetDepositsRequest{Id: &id, Page: &page})

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, depositsStatus)
	return
}

func (controller *depositsController) NewV1Deposits(context echo.Context) (err error) {
	id := context.Get("userId").(int64)
	newDepositRequest := &protoSchema.CreateDepositRequest{}

	if err = context.Bind(newDepositRequest); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return err
	}

	newDepositRequest.Id = &id

	newDepositStatus, err := controller.microservice.CreateDeposit(microContext.TODO(),
		newDepositRequest)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, newDepositStatus)
	return
}
