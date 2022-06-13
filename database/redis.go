package database

import (
	"github.com/gofiber/fiber/v2"
	fr "github.com/gofiber/storage/redis"
)

var FiberStorage fiber.Storage

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
