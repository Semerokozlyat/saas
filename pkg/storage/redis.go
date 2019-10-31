package storage

import (
	"fmt"
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

func (storage *RedisStorage) Get(id string) ([]byte, error) {
	conn := storage.pool.Get()
	defer conn.Close()

	selectRes, err := redis.String(conn.Do("SELECT", storage.db))
	if selectRes != "OK" {
		return nil, fmt.Errorf("failed to select Redis DB on GET: %v", err)
	}

	data, err := redis.Bytes(conn.Do("GET", id))
	if err != nil {
		return nil, fmt.Errorf("failed to get data from db")
	}

	return data, nil
}

func (storage *RedisStorage) Put(id string, data []byte) error {
	conn := storage.pool.Get()
	defer conn.Close()

	selectRes, err := redis.String(conn.Do("SELECT", storage.db))
	if selectRes != "OK" {
		return fmt.Errorf("failed to select Redis DB on SET: %v", err)
	}

	exists, err := redis.Int(conn.Do("EXISTS", id))
	if err != nil {
		return fmt.Errorf("failed to check if DB key already exists")
	} else if exists == 1 {
		return fmt.Errorf("DB key already exists")
	}
	result, err := redis.String(conn.Do("SET", id, data))
	if err != nil || result != "OK" {
		return fmt.Errorf("failed to put data into DB")
	}

	return nil
}

func (storage *RedisStorage) Update(id string) error {
	return nil
}

func (storage *RedisStorage) Delete(id string) (int, error) {
	conn := storage.pool.Get()
	defer conn.Close()

	selectRes, err := redis.String(conn.Do("SELECT", storage.db))
	if selectRes != "OK" {
		return 0, fmt.Errorf("failed to select Redis DB on DELETE: %v", err)
	}
	delCount, err := redis.Int(conn.Do("DEL", id))
	if err != nil {
		return 0, fmt.Errorf("failed to delete data from DB: %v", err)
	}
	return delCount, nil
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
		addr:    addrVal,
		db:      dbVal,
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
