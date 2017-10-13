package handlers

import (
	"auth/services"
	"bytes"
	"crypto/sha256"
	"errors"
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/net/context"
	"helpers/values"
	proto "protos"
	"strconv"
	"time"
)

const tokensLength int = 32

type authHandler struct {
	proto.AuthHandler
	storage services.AuthStorageService
}

// Factory that creates new AuthHandler instance
func NewAuthHandler(storage services.AuthStorageService) proto.AuthHandler {
	return &authHandler{storage: storage}
}

// -----------MAIN MICROSERVICE METHODS------------

// AuthHandler Login method implementation
func (handler *authHandler) Login(ctx context.Context, req *proto.AuthLoginRequest, rsp *proto.AuthTokenData) (err error) {
	userName := req.GetUsername()
	password := req.GetPassword()

	user, err := handler.getUser(&userName)

	if err != nil {
		return err
	}

	passwordSum := handler.generatePasswordSum(password, user.Salt)

	if bytes.Compare(passwordSum, user.Password) != 0 {
		return errors.New("password doesn't match")
	}

	// invalidate other sessions
	err = handler.invalidateUserSessions(user.Id)

	if err != nil {
		return err
	}

	rsp.Token = handler.generateSession(user.Salt, user.Id)

	return nil
}

func (handler *authHandler) Logout(ctx context.Context, req *proto.AuthLogoutRequest, rsp *proto.AuthLogoutResponse) error {
	userId := req.GetId()
	err := handler.invalidateUserSessions(userId)

	if err != nil {
		return err
	}

	return nil
}

func (handler *authHandler) Register(ctx context.Context, req *proto.AuthRegisterRequest, rsp *proto.AuthTokenData) (err error) {
	var referrerId int64

	password := req.GetPassword()
	leg := req.GetLeg()
	referrer := req.GetReferrer()

	if len(leg) != 0 && len(referrer) != 0 {
		referrerId, err = handler.validateReferral(leg, referrer)

		if referrerId == 0 {
			return err
		}
	}

	salt := handler.generateSalt()
	passwordSum := handler.generatePasswordSum(password, salt)
	var userId int64

	timeNow := time.Now().UnixNano()

	username := req.GetUsername()

	if handler.validateUsername(username) == false {
		return errors.New("Invalid username")
	}

	err = handler.storage.QueryGet(&userId, "RegisterUser", username, req.GetPhone(), req.GetEmail(), salt, passwordSum, timeNow)

	if err != nil {
		return err
	}

	rsp.Token = handler.generateSession(salt, userId)

	go handler.storage.Exec("CreateWallets", userId)
	go handler.registerReferrer(userId, referrerId, leg)

	return nil
}

func (handler *authHandler) ProlongSession(ctx context.Context, req *proto.AuthTokenData, rsp *proto.AuthProlongSessionResponse) error {
	token := req.GetToken()

	id, err := handler.storage.GetUserId(token)

	if err != nil {
		return err
	}

	rsp.Id = &id

	return nil
}

// -----------HELPER METHODS------------

func (handler *authHandler) generatePasswordSum(password string, salt []byte) []byte {
	passwordSum, _ := scrypt.Key([]byte(password), salt, 16384, 8, 1, tokensLength)
	return passwordSum
}

func (handler *authHandler) validateReferral(leg string, referrer string) (int64, error) {
	if _, ok := values.Legs["leg"]; ok != false {
		return 0, errors.New("Unknown leg")
	}

	var id int64

	err := handler.storage.QueryGet(&id, "ValidateReferrer", referrer)

	return id, err
}

func (handler *authHandler) generateSalt() []byte {
	timeOfRegistration, _ := time.Now().MarshalBinary()
	salt := sha256.Sum256(timeOfRegistration)
	return salt[:]
}

func (handler *authHandler) generateSession(salt []byte, id int64) *string {
	stringifiedToken, validTill := handler.storage.GenerateSessionToken(salt)

	// save to db
	go handler.storage.Exec("SaveSessionToken", id, stringifiedToken, validTill)

	// cache in redis. Redis is bad with concurrency
	handler.storage.CacheToken(id, stringifiedToken, validTill)

	// return string
	return &stringifiedToken
}

func (handler *authHandler) invalidateUserSessions(id int64) (err error) {
	timeNow := time.Now().UnixNano()
	stringTimeNow := strconv.FormatInt(timeNow, 10)

	go handler.storage.Exec("InvalidateOtherSessions", id, stringTimeNow)

	return handler.storage.CacheSessionInvalidate(id)
}

func (handler *authHandler) getUser(username *string) (*struct {
	Id       int64
	Salt     []byte
	Password []byte
}, error) {
	userStructure := struct {
		Id       int64
		Salt     []byte
		Password []byte
	}{}

	err := handler.storage.QueryGet(&userStructure, "GetUserCredentials", *username)

	return &userStructure, err
}

func (handler *authHandler) registerReferrer(userId int64, referrer int64, leg string) error {
	var legValue interface{}
	var referrerValue interface{}

	if leg != "" {
		legValue = leg
	}

	if referrer != 0 {
		referrerValue = referrer
	}

	return handler.storage.Exec("RegisterReferrer", userId, referrerValue, legValue)
}

func (handler *authHandler) validateUsername(username string) bool {
	return govalidator.Matches(username, "^[a-zA-Z0-9_-]{3,15}$")
}
