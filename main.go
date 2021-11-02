package main

import (
	"time"
	"userland/api"
	"userland/store"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	// TODO use external config management (toml?)
	_ = godotenv.Load(".env")

	serverCfg := api.ServerConfig{
		Host:            "0.0.0.0",
		Port:            "80",
		ReadTimeout:     500 * time.Millisecond,
		WriteTimeout:    500 * time.Millisecond,
		ShutdownTimeout: 10 * time.Second,
	}
	postgresCfg := store.PostgresConfig{
		Host:     "db_userland",
		Port:     5432,
		Username: "admin",
		Password: "admin",
		Database: "userland",
	}
	postgresDB, err := store.NewPG(postgresCfg)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}

	redisCfg := store.RedisConfig{
		Host:     "redis",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	redisDb, err := store.NewRedis(redisCfg)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}

	serverDataSource := &api.DataSource{
		PostgresDB: postgresDB,
		RedisDB:    redisDb,
	}

	srv := api.NewServer(serverCfg, serverDataSource)
	srv.Start()
}
