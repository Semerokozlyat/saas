package storage

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisStorage struct {
	network string
	host    string
	port    string
	db      string

	pool *redis.Pool
}

func (storage *RedisStorage) Get(id string) (interface{}, error) {
	return nil, nil
}

func (storage *RedisStorage) Put(interface{}) error {
	return nil
}

func (storage *RedisStorage) Update(id string) error {
	return nil
}

func (storage *RedisStorage) Delete(id string) error {
	return nil
}

func (storage *RedisStorage) Create() *RedisStorage {
	resultedStorage := &RedisStorage{
		host: "127.0.0.1",
		port: "6379",
		db:   "1",
	}
	storage.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(storage.network, storage.host)
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
	return resultedStorage
}
