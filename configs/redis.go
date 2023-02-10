package configs

import (
	"os"
	"runtime"
	"strconv"

	"github.com/gofiber/storage/redis"
)

var RedisStorage *redis.Storage

func (c *Config) RedisConfig() {
	var (
		host     string = os.Getenv("REDIS_HOST")
		port     string = os.Getenv("REDIS_PORT")
		user     string = os.Getenv("REDIS_USER")
		password string = os.Getenv("REDIS_PASSWORD")
	)

	redisPort, _ := strconv.Atoi(port)
	redisConfig := redis.Config{
		Host:      host,
		Port:      redisPort,
		Username:  user,
		Password:  password,
		URL:       "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	}

	RedisStorage = redis.New(redisConfig)
}
