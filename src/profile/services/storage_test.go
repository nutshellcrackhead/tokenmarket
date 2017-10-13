package services_test

import (
	"log"
	"profile/services"
	"testing"
	"tests/mocks/helpers"
	"tests/mocks/profile"
)

func TestProfileStorageService(t *testing.T) {
	// arrange
	db := new(helpers.DatabaseMock)
	queries := map[string]string{
		"GetProfile":        "GetProfileQuery",
		"UpdateProfileData": "UpdateProfileDataQuery",
	}
	logger := &log.Logger{}

	t.Run("should be able to close connection to database", func(t *testing.T) {
		testProfileStorage := services.NewProfileStorageService(db, logger, queries)
		db.On("Close").Return()

		// act
		testProfileStorage.Close()

		// assert
		db.AssertCalled(t, "Close")
	})

	t.Run("should be able to pass retrieved request to database helper", func(t *testing.T) {
		testProfileStorage := services.NewProfileStorageService(db, logger, queries)

		profileInstance := &profile.ProfileMock{}
		params := []interface{}{1}

		mockedArguments := append([]interface{}{
			profileInstance,
			"GetProfileQuery",
		}, params...)

		db.On("GetRecord", mockedArguments...).Return()

		// act
		testProfileStorage.QueryGet(profileInstance, "GetProfile", params...)

		// assert
		db.AssertCalled(t, "GetRecord", mockedArguments...)
	})
}
