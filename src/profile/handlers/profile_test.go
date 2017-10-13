package handlers_test

import (
	"github.com/stretchr/testify/assert"
	"profile/handlers"
	"profile/models"
	protos "protos"
	"testing"
	profileMocks "tests/mocks/profile"
	"tests/mocks/vendor"
)

func TestProfileHandler(t *testing.T) {
	profileMock := &profileMocks.ProfileMock{}
	profileMock.On("GetId").Return(int32(123))
	profileMock.On("GetName").Return("Name")
	profileMock.On("GetLocation").Return("kiev")
	profileMock.On("GetAvatar").Return("avatar")
	profileMock.On("GetEmail").Return("test@test.com")
	profileMock.On("GetSkype").Return("skype")
	profileMock.On("GetPhone").Return("phonenumber")

	profileInstanceMockFactory := func(id int32) models.Profile {
		return profileMock
	}

	t.Run("should be able to get the profile data and send it back", func(t *testing.T) {
		// arrange
		req := &protos.ProfileGetRequest{}
		rsp := &protos.ProfileData{}

		profileHandler := handlers.NewProfileHandler(profileInstanceMockFactory)

		// act
		profileHandler.Get(&vendor.Context{}, req, rsp)

		// assert
		assert.Equal(t, rsp.GetId(), int32(123), "Id should be equal")
		assert.Equal(t, rsp.GetName(), "Name", "Name should be equal")
		assert.Equal(t, rsp.GetLocation(), "kiev", "Location should be equal")
		assert.Equal(t, rsp.GetAvatar(), "avatar", "Avatar should be equal")
		assert.Equal(t, rsp.GetEmail(), "test@test.com", "Email should be equal")
		assert.Equal(t, rsp.GetSkype(), "skype", "Skype should be equal")
		assert.Equal(t, rsp.GetPhone(), "phonenumber", "Phone number should be equal")
	})
}
