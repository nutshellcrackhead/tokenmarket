package models

import (
	"database/sql"
	"log"
	"profile/services"
	proto "protos"
	"reflect"
	"regexp"
)

type Profile interface {
	Fetch() error
	Save() error
	GetCityList(name string, citiesList *[]*proto.City) error
	GetCity(city *proto.City) error
	GetId() int64
	GetLocation() *proto.City
	GetUsername() *string
	GetAvatar() *string
	GetPhone() *string
	GetSkype() *string
	GetEmail() *string
	GetMuted() *bool
	SetLocation(value *int64) error
	SetUsername(value *string)
	SetAvatar(value *string)
	SetPhone(value *string)
	SetSkype(value *string)
	SetMuted(value *bool)
}

type ProfileData interface {
	GetProfile(p Profile) error
	SaveProfile(p Profile) error
	GetCityList(name string, citiesList *[]*proto.City) error
	GetCity(city *proto.City) error
}

type profileData struct {
	storage services.ProfileStorageService
	logger  log.Logger
}

type profile struct {
	Id                          int64
	Location                    *proto.City `db:"token_locations"`
	Username                    string
	Avatar, Email, Phone, Skype sql.NullString
	Salt, Password              []byte
	profileData                 ProfileData
	queries                     map[string]string
	Muted                       bool
	alphaRegex                  *regexp.Regexp
}

func NewProfile(profileData ProfileData, id int64) Profile {
	profileItem := &profile{
		Id:          id,
		profileData: profileData,
		queries:     services.GetQueries(),
	}

	return profileItem
}

func NewProfileData(storage services.ProfileStorageService) ProfileData {
	return &profileData{storage, log.Logger{}}
}

// Sync profile model
func (p *profile) Fetch() error {
	return p.profileData.GetProfile(p)
}

func (p *profile) Save() error {
	return p.profileData.SaveProfile(p)
}

func (p *profile) GetId() int64 {
	return p.Id
}

func (p *profile) GetMuted() *bool {
	return &p.Muted
}

func (p *profile) GetLocation() *proto.City {
	return p.Location
}

func (p *profile) GetUsername() *string {
	return &p.Username
}

func (p *profile) GetAvatar() *string {
	return p.getNullStringValue("Avatar")
}

func (p *profile) GetPhone() *string {
	return p.getNullStringValue("Phone")
}

func (p *profile) GetSkype() *string {
	return p.getNullStringValue("Skype")
}

func (p *profile) GetEmail() *string {
	return p.getNullStringValue("Email")
}

func (p *profile) getNullStringValue(valueName string) *string {
	profileReflect := reflect.ValueOf(p).Elem()
	valueToGet := profileReflect.FieldByName(valueName)

	if !valueToGet.IsValid() {
		return nil
	}

	valid := valueToGet.FieldByName("Valid").Bool()

	if !valid {
		return nil
	}

	value := valueToGet.FieldByName("String").String()

	return &value
}

func (p *profile) SetLocation(value *int64) error {
	if value == nil {
		return nil
	}

	if p.Location == nil {
		p.Location = &proto.City{}
	}

	p.Location.Id = value
	return p.profileData.GetCity(p.Location)
}

func (p *profile) SetUsername(value *string) {
	p.Username = *value
}

func (p *profile) SetAvatar(value *string) {
	p.Avatar = p.setNullStringValue(value)
}

func (p *profile) SetPhone(value *string) {
	p.Phone = p.setNullStringValue(value)
}

func (p *profile) SetSkype(value *string) {
	p.Skype = p.setNullStringValue(value)
}

func (p *profile) SetEmail(value *string) {
	p.Email = p.setNullStringValue(value)
}

func (p *profile) SetMuted(value *bool) {
	p.Muted = *value
}

func (p *profile) setNullStringValue(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{Valid: true, String: *value}
}

func (p *profile) setNullIntValue(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{Valid: false}
	}

	return sql.NullInt64{Valid: true, Int64: *value}
}

func (p *profile) GetCityList(name string, citiesList *[]*proto.City) error {
	return p.profileData.GetCityList(name, citiesList)
}

func (p *profile) GetCity(city *proto.City) error {
	return p.profileData.GetCity(city)
}

// GetProfile data from storage
func (data *profileData) GetProfile(p Profile) error {
	return data.storage.QueryGet(p, "GetProfile", p.GetId())
}

func (data *profileData) GetCityList(name string, citiesList *[]*proto.City) error {
	return data.storage.QuerySelect(citiesList, "GetCitiesList", name)
}

// SaveProfile data to storage
func (data *profileData) SaveProfile(p Profile) error {
	location := p.GetLocation()
	var locationIdReference *int64

	if location != nil {
		locationId := location.GetId()

		if locationId != 0 {
			locationIdReference = &locationId
		}
	}

	return data.storage.QueryGet(p, "UpdateProfileData", p.GetId(), locationIdReference,
		p.GetSkype(), p.GetAvatar())
}

func (data *profileData) GetCity(city *proto.City) error {
	return data.storage.QueryGet(city, "GetCity", city.GetId())
}
