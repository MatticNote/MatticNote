package config

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis"
	redigo "github.com/gomodule/redigo/redis"
)

func GetFiberRedisMemory() fiber.Storage {
	return redis.New(redis.Config{
		Host:     Config.Redis.Address,
		Port:     int(Config.Redis.Port),
		Username: Config.Redis.Username,
		Password: Config.Redis.Password,
		Database: Config.Redis.Database,
		Reset:    false,
	})
}

func GetRedisDial() (redigo.Conn, error) {
	return redigo.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", Config.Redis.Address, Config.Redis.Port),
		redigo.DialDatabase(Config.Redis.Database),
		redigo.DialUsername(Config.Redis.Username),
		redigo.DialPassword(Config.Redis.Password),
	)
}
