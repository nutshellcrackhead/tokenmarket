package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	microContext "golang.org/x/net/context"
	"net/http"
	protoSchema "protos"
	"unicode/utf8"
)

type ProfileServerController interface {
	GetV1ProfileMy(context echo.Context) error
	PatchV1ProfileMy(context echo.Context) error
	CitiesV1Profile(context echo.Context) error
}

type profileController struct {
	microservice protoSchema.ProfileClient
}

const cityNameMinLength int = 2

func NewProfileController(microservice protoSchema.ProfileClient) ProfileServerController {
	return &profileController{microservice: microservice}
}

func (controller *profileController) GetV1ProfileMy(context echo.Context) (err error) {
	id := context.Get("userId").(int64)
	profile, err := controller.microservice.Get(microContext.TODO(), &protoSchema.ProfileGetRequest{Id: &id})

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, profile)
	return
}

func (controller *profileController) PatchV1ProfileMy(context echo.Context) (err error) {
	id := context.Get("userId").(int64)
	patchRequest := &protoSchema.ProfilePatchRequest{}

	if err = context.Bind(&patchRequest); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusForbidden, map[string]string{"error": "invalid_data"})
		return err
	}

	patchRequest.Id = &id

	profile, err := controller.microservice.Patch(microContext.TODO(), patchRequest)

	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	context.JSON(http.StatusOK, profile)
	return
}

func (controller *profileController) CitiesV1Profile(context echo.Context) (err error) {
	cityFragment := context.QueryParam("name")

	if cityFragment == "" || utf8.RuneCountInString(cityFragment) < cityNameMinLength {
		return context.NoContent(http.StatusNoContent)
	}

	citiesRequest := &protoSchema.ProfileCitiesRequest{CityFragment: &cityFragment}
	citiesList, err := controller.microservice.Cities(microContext.TODO(), citiesRequest)

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "unprocessable_entity"})
		return
	}

	if len(citiesList.GetCities()) == 0 {
		context.NoContent(http.StatusNoContent)
		return
	}

	context.JSON(http.StatusOK, citiesList)
	return
}
