package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	microContext "golang.org/x/net/context"
	"net/http"
	protoSchema "protos"
	"strconv"
)

type ReferralServerController interface {
	GetV1ReferralMy(context echo.Context) error
}

type referralController struct {
	microservice protoSchema.ReferralClient
}

func NewReferralController(microservice protoSchema.ReferralClient) ReferralServerController {
	return &referralController{microservice: microservice}
}

func (controller *referralController) GetV1ReferralMy(context echo.Context) (err error) {
	page := int64(0)

	id := context.Get("userId").(int64)
	pageString := context.QueryParam("page")
	parsedPage, parsePageErr := strconv.ParseInt(pageString, 10, 64)

	if parsePageErr == nil {
		page = parsedPage
	}

	referralStatus, err := controller.microservice.List(microContext.TODO(),
		&protoSchema.ReferralListRequest{Id: &id, Page: &page})

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, referralStatus)
	return
}
