package database

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"time"
)

type Cache interface {
	Set(key string, value interface{}) (reply interface{}, err error)
	Get(key string) (reply interface{}, err error)
	Delete(key string) (reply interface{}, err error)
	HashmapSet(key string, valueMap map[string]interface{}) (reply interface{}, err error)
	GetBySelector(selector string) ([]interface{}, error)
	Close()
}

type redisCache struct {
	connection *redis.Pool
}

func getCacheConfigs() *struct {
	Network  string
	Address  string
	Password string
} {
	network := viper.GetString("cache.network")
	address := viper.GetString("cache.address")
	password := viper.GetString("cache.password")

	cacheConfigs := struct {
		Network  string
		Address  string
		Password string
	}{network, address, password}

	return &cacheConfigs
}

// open new connect to Postgres database
func NewCache() Cache {
	initEnvConfigs()
	cacheConfigs := getCacheConfigs()

	cachePassword := redis.DialPassword(cacheConfigs.Password)
	connectionTimeout := redis.DialConnectTimeout(0)

	connection := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(cacheConfigs.Network, cacheConfigs.Address, cachePassword, connectionTimeout)
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &redisCache{connection}
}

func (cache *redisCache) GetBySelector(selector string) ([]interface{}, error) {
	connection, _ := cache.connection.Dial()
	values, err := redis.Values(connection.Do("SCAN", 0, "MATCH", selector, "COUNT", 1000))

	connection.Close()
	if len(values) < 2 {
		return []interface{}{}, nil
	}

	return redis.Values(values[1], err)
}

func (cache *redisCache) Set(key string, value interface{}) (reply interface{}, err error) {
	connection, _ := cache.connection.Dial()
	reply, err = connection.Do("SET", key, value)
	connection.Close()

	if err != nil {
		fmt.Println("Set", key, value, err)
	}

	return
}

func (cache *redisCache) Get(key string) (reply interface{}, err error) {
	connection, _ := cache.connection.Dial()
	reply, err = connection.Do("GET", key)
	connection.Close()

	if err != nil {
		fmt.Println("Get", key, err)
	}

	return
}

func (cache *redisCache) HashmapSet(key string, valueMap map[string]interface{}) (reply interface{}, err error) {
	connection, _ := cache.connection.Dial()
	reply, err = connection.Do("HMSET", redis.Args{key}.AddFlat(valueMap)...)
	connection.Close()

	if err != nil {
		fmt.Println("HMSET", key, valueMap)
	}

	return
}

func (cache *redisCache) Close() {
}

func (cache *redisCache) Delete(key string) (interface{}, error) {
	connection, _ := cache.connection.Dial()
	result, err := connection.Do("DEL", key)
	connection.Close()

	if err != nil {
		fmt.Println("DEL", key)
	}

	return result, err
}
