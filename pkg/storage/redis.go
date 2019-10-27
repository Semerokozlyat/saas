package storage

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"time"
)

type RedisStorage struct {
	network string
	addr    string
	db      string

	pool *redis.Pool
}

func (storage *RedisStorage) Get(id string) (interface{}, error) {
	return nil, nil
}

func (storage *RedisStorage) Put(interface{}) error {
	log.Print("Storage Put operation")
	return nil
}

func (storage *RedisStorage) Update(id string) error {
	return nil
}

func (storage *RedisStorage) Delete(id string) error {
	return nil
}

func CreateRedisStorage() *RedisStorage {

	log.Print("Creating Redis storage")

	networkVal, found := os.LookupEnv("REDIS_NETWORK")
	if !found {
		networkVal = "tcp"
	}

	addrVal, found := os.LookupEnv("REDIS_ADDR")
	if !found {
		addrVal = "redis:6379"
	}

	dbVal, found := os.LookupEnv("REDIS_DB")
	if !found {
		dbVal = "2"
	}

	resultedStorage := &RedisStorage{
		network: networkVal,
		addr: addrVal,
		db:   dbVal,
	}

	log.Printf("Redis storage connection: %s/%s db %s", networkVal, addrVal, dbVal)

	resultedStorage.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(resultedStorage.network, resultedStorage.addr)
		},
		DialContext: nil,
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
		MaxIdle:         20,
		MaxActive:       0,
		IdleTimeout:     20 * time.Second,
		Wait:            false,
		MaxConnLifetime: 60 * time.Second,
	}
	log.Print("Redis storage created")
	return resultedStorage
}
