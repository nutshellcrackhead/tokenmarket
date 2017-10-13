package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	microContext "golang.org/x/net/context"
	"net/http"
	protoSchema "protos"
)

type OperationsServerController interface {
	NewV1Operations(context echo.Context) error
	StatusV1Operations(context echo.Context) error
	ActivateV1Operations(context echo.Context) error
	NewV1PayoutOperations(context echo.Context) error
}

type operationsController struct {
	microservice protoSchema.OperationsClient
}

func NewOperationsController(microservice protoSchema.OperationsClient) OperationsServerController {
	return &operationsController{microservice: microservice}
}

func (controller *operationsController) NewV1Operations(context echo.Context) (err error) {
	id := context.Get("userId").(int64)

	createPaymentRequest := protoSchema.CreatePaymentRequest{}

	if err = context.Bind(&createPaymentRequest); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusForbidden, map[string]string{"error": "invalid_data"})
		return err
	}

	createPaymentRequest.Id = &id

	operationsStatus, err := controller.microservice.CreatePayment(microContext.TODO(),
		&createPaymentRequest)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, operationsStatus)
	return
}

func (controller *operationsController) NewV1PayoutOperations(context echo.Context) (err error) {
	id := context.Get("userId").(int64)

	createPayoutRequest := protoSchema.CreatePayoutRequest{}

	if err = context.Bind(&createPayoutRequest); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusForbidden, map[string]string{"error": "invalid_data"})
		return err
	}

	createPayoutRequest.Id = &id

	operationsStatus, err := controller.microservice.CreatePayout(microContext.TODO(),
		&createPayoutRequest)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, operationsStatus)
	return
}

func (controller *operationsController) StatusV1Operations(context echo.Context) (err error) {
	updatePaymentRequest := protoSchema.UpdatePaymentRequest{}

	if err = context.Bind(&updatePaymentRequest); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusForbidden, map[string]string{"error": "invalid_data"})
		return err
	}

	operationsStatus, err := controller.microservice.UpdatePaymentStatus(microContext.TODO(),
		&updatePaymentRequest)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, operationsStatus)
	return
}

func (controller *operationsController) ActivateV1Operations(context echo.Context) (err error) {
	id := context.Get("userId").(int64)

	createPaymentRequest := protoSchema.CreatePaymentRequest{}

	if err = context.Bind(&createPaymentRequest); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusForbidden, map[string]string{"error": "invalid_data"})
		return err
	}

	createPaymentRequest.Id = &id

	operationsStatus, err := controller.microservice.ActivateUser(microContext.TODO(),
		&createPaymentRequest)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, operationsStatus)
	return
}
