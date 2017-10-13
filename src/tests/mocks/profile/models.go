package profile

import (
	"github.com/stretchr/testify/mock"
	"profile/models"
)

/**
PROFILE MOCK
*/
type ProfileMock struct {
	mock.Mock
	Id                                          int32
	Queries                                     map[string]string
	Avatar, Email, Location, Phone, Skype, Name string
}

func (p *ProfileMock) Fetch() {
	p.Called()
}

func (p *ProfileMock) GetId() int32 {
	args := p.Called()

	return args.Get(0).(int32)
}

func (p *ProfileMock) GetLocation() string {
	args := p.Called()

	return args.String(0)
}

func (p *ProfileMock) GetName() string {
	args := p.Called()

	return args.String(0)
}

func (p *ProfileMock) GetAvatar() string {
	args := p.Called()

	return args.String(0)
}

func (p *ProfileMock) GetPhone() string {
	args := p.Called()

	return args.String(0)
}

func (p *ProfileMock) GetSkype() string {
	args := p.Called()

	return args.String(0)
}

func (p *ProfileMock) GetEmail() string {
	args := p.Called()

	return args.String(0)
}

/**
PROFILE DATA MOCK
*/
type ProfileDataMock struct {
	mock.Mock
}

func (data *ProfileDataMock) GetProfile(p models.Profile) {
	data.Called(p)
}
