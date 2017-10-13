package handlers

import (
	"golang.org/x/net/context"
	"profile/models"
	proto "protos"
)

type profileHandler struct {
	proto.ProfileHandler
	profileConstructor func(id int64) models.Profile
}

// Factory that creates new ProfileHandler instance
func NewProfileHandler(profileConstructor func(id int64) models.Profile) proto.ProfileHandler {
	return &profileHandler{profileConstructor: profileConstructor}
}

// Retrieves profile instance by id
func (handler *profileHandler) getProfileInstance(id int64) models.Profile {
	return handler.profileConstructor(id)
}

// ProfileHandler Get method implementation
func (handler *profileHandler) Get(ctx context.Context, req *proto.ProfileGetRequest, rsp *proto.ProfileData) (err error) {
	profileInstance := handler.getProfileInstance(req.GetId())
	err = profileInstance.Fetch()

	handler.setProfileDataToResponse(profileInstance, rsp)

	return err
}

// ProfileHandler Get method implementation
func (handler *profileHandler) Patch(ctx context.Context, req *proto.ProfilePatchRequest, rsp *proto.ProfileData) (err error) {
	profileInstance := handler.getProfileInstance(req.GetId())
	profileInstance.SetAvatar(req.Avatar)
	err = profileInstance.SetLocation(req.Location)

	if err != nil {
		return err
	}

	profileInstance.SetSkype(req.Skype)
	err = profileInstance.Save()

	if err != nil {
		return err
	}

	handler.setProfileDataToResponse(profileInstance, rsp)

	return nil
}

func (handler *profileHandler) Cities(ctx context.Context, req *proto.ProfileCitiesRequest, rsp *proto.ProfileCitiesResponse) error {
	cityFragmentName := req.GetCityFragment()
	handler.setCitiesNameList(cityFragmentName, rsp)

	return nil
}

func (handler *profileHandler) setCitiesNameList(cityFragment string, rsp *proto.ProfileCitiesResponse) {
	cities := []*proto.City{}
	handler.getProfileInstance(int64(0)).GetCityList(cityFragment, &cities)

	rsp.Cities = cities
}

func (handler *profileHandler) setProfileDataToResponse(profile models.Profile, response *proto.ProfileData) {
	id := profile.GetId()
	location := profile.GetLocation()

	response.Id = &id
	response.Avatar = profile.GetAvatar()
	response.Email = profile.GetEmail()
	response.Location = location
	response.Phone = profile.GetPhone()
	response.Skype = profile.GetSkype()
	response.Username = profile.GetUsername()
	response.Muted = profile.GetMuted()
}
