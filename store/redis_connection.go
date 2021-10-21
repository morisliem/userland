package store

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Host     string
	DB       int
	Port     int
	Password string
}

func NewRedis(config RedisConfig) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return client, err
	}

	return client, nil
}
