package database

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	fr "github.com/gofiber/storage/redis"
	"github.com/gomodule/redigo/redis"
)

var (
	FiberStorage fiber.Storage
	RedisPool    *redis.Pool
)

func InitFiberRedisMemory(
	host string,
	port uint16,
	username string,
	password string,
	database int,
) {
	FiberStorage = fr.New(fr.Config{
		Host:     host,
		Port:     int(port),
		Username: username,
		Password: password,
		Database: database,
		Reset:    false,
	})
}

func InitRedis(
	host string,
	port uint16,
	username string,
	password string,
	database int,
) {
	RedisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", host, port),
				redis.DialUsername(username),
				redis.DialPassword(password),
				redis.DialDatabase(database),
			)
		},
		MaxIdle:         16,
		MaxActive:       0,
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}
