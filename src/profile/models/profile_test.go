package models_test

import (
	"github.com/stretchr/testify/assert"
	"profile/models"
	"testing"
	profileMock "tests/mocks/profile"
)

func TestProfileDataModel(t *testing.T) {
	storage := &profileMock.ProfileStorageServiceMock{}
	profile := &profileMock.ProfileMock{}

	t.Run("should be able to get profile", func(t *testing.T) {
		// arrange
		userId := int32(1)

		storage.On("QueryGet", profile, "GetProfile", userId).Return()
		profile.On("GetId").Return(userId)

		// act
		profileData := models.NewProfileData(storage)
		profileData.GetProfile(profile)

		// assert
		storage.AssertCalled(t, "QueryGet", profile, "GetProfile", userId)
	})
}

func TestProfileModel(t *testing.T) {
	profileDataMock := &profileMock.ProfileDataMock{}

	profileId := int32(1023121)
	profileInstance := models.NewProfile(profileDataMock, profileId)

	t.Run("should be able to return id", func(t *testing.T) {
		// assert
		assert.Equal(t, profileInstance.GetId(), profileId, "should be equal")
	})

	t.Run("should be able to fetch data from database", func(t *testing.T) {
		// act
		profileDataMock.On("GetProfile", profileInstance).Return()
		profileInstance.Fetch()

		// assert
		profileDataMock.AssertCalled(t, "GetProfile", profileInstance)
	})

	t.Run("should be able to return id", func(t *testing.T) {
		// act
		profileDataMock.On("GetProfile", profileInstance).Return()
		profileInstance.Fetch()

		// assert
		profileDataMock.AssertCalled(t, "GetProfile", profileInstance)
	})
}
