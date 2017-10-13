package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/garyburd/redigo/redis"
	"helpers/database"
	"html/template"
	"log"
	"strconv"
	"strings"
	"time"
)

type authStorageService struct {
	db                database.Database
	cache             database.Cache
	logger            *log.Logger
	queries           map[string]string
	templates         map[string]*template.Template
	sessionTimeToLive time.Duration
}

type AuthStorageService interface {
	Close()
	QueryGet(p interface{}, queryName string, params ...interface{}) error
	CacheToken(id int64, token string, validTill int64) error
	CacheSessionInvalidate(id int64) error
	Exec(queryName string, params ...interface{}) error
	GetUserId(token string) (int64, error)
	GenerateSessionToken(salt []byte) (string, int64)
}

// Profile Storage Service
func NewAuthStorageService(db database.Database, cache database.Cache, logger *log.Logger, queries map[string]string, sessionTimeToLive time.Duration) AuthStorageService {
	templates := make(map[string]*template.Template)
	templates["sessionKey"], _ = template.New("session").Parse(sessionKeyTemplate)

	authStorage := &authStorageService{db, cache, logger, queries, templates, sessionTimeToLive}

	return authStorage
}

// Close the database connection
func (storage *authStorageService) Close() {
	storage.db.Close()
	storage.cache.Close()
}

func (storage *authStorageService) QueryGet(p interface{}, queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.GetRecord(p, query, params...)
}

func (storage *authStorageService) Exec(queryName string, params ...interface{}) error {
	query := storage.queries[queryName]
	return storage.db.Exec(query, params...)
}

const sessionKeyTemplate string = `session:{{ .token }}:{{ .id }}`

func (storage *authStorageService) CacheToken(id int64, token string, validTill int64) error {
	storage.CacheSessionInvalidate(id)

	session := map[string]interface{}{
		"id":    id,
		"token": token,
	}

	key := bytes.Buffer{}
	storage.templates["sessionKey"].Execute(&key, session)

	_, err := storage.cache.Set(key.String(), validTill)

	return err
}

func (storage *authStorageService) CacheSessionInvalidate(id int64) error {
	session := map[string]interface{}{
		"id":    id,
		"token": "*",
	}

	selector := bytes.Buffer{}
	storage.templates["sessionKey"].Execute(&selector, session)

	values, err := storage.cache.GetBySelector(selector.String())

	if err != nil {
		return err
	}

	for _, value := range values {
		convertedValue := value.([]byte)
		key := string(convertedValue[:])
		storage.cache.Delete(key)
	}

	return err
}

func (storage *authStorageService) GetUserId(token string) (int64, error) {
	var id int64

	selector, substring := storage.getCacheSubstringAndSelector(token)
	keys, err := storage.cache.GetBySelector(selector)

	if err != nil || len(keys) == 0 {
		return id, err
	}

	tokenKey := keys[0]

	timeNow := time.Now().UnixNano()
	key := storage.getCacheKey(tokenKey)
	id = storage.getUserIdFromCacheKey(key, substring)
	validTill, err := storage.cache.Get(key)
	validTillInt, err := redis.Int64(validTill, err)

	if err != nil {
		return id, err
	}

	if timeNow > validTillInt {
		storage.cache.Delete(key)
		return id, errors.New("Expired token")
	}

	newTime := time.Now().Add(10 * time.Minute).UnixNano()
	_, err = storage.cache.Set(key, newTime)
	go storage.db.Exec(storage.queries["ProlongSession"], token, newTime)

	return id, err
}

func (storage *authStorageService) getCacheSubstringAndSelector(token string) (selector string, substring string) {
	session := map[string]interface{}{
		"id":    "",
		"token": token,
	}

	substringBuffer := bytes.Buffer{}
	storage.templates["sessionKey"].Execute(&substringBuffer, session)

	substring = substringBuffer.String()

	selectorBuffer := bytes.Buffer{}
	session["id"] = "*"
	storage.templates["sessionKey"].Execute(&selectorBuffer, session)

	return selectorBuffer.String(), substring
}

func (storage *authStorageService) getCacheKey(tokenKey interface{}) string {
	convertedKey := tokenKey.([]byte)
	return string(convertedKey[:])
}

func (storage *authStorageService) getUserIdFromCacheKey(key string, substring string) int64 {
	stringifiedId := strings.TrimPrefix(key, substring)
	id, _ := strconv.ParseInt(stringifiedId, 10, 64)
	return id
}

func (storage *authStorageService) GenerateSessionToken(salt []byte) (string, int64) {
	// generate token
	validTill := time.Now().Add(time.Minute * storage.sessionTimeToLive).UnixNano()
	stringifiedValidTill := strconv.FormatInt(validTill, 10)

	sessionToken := sha256.New()
	sessionToken.Write([]byte(stringifiedValidTill))
	sessionToken.Write(salt)
	shaToken := sessionToken.Sum(nil)

	return base64.StdEncoding.EncodeToString(shaToken[:]), validTill
}
