package main

import (
	"log"
	"time"
	"userland/api"
	"userland/store"

	"github.com/joho/godotenv"
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
		// TODO proper logging with zlogger
		log.Fatalf("failed to open db conn: %v\n", err)
	}

	redisCfg := store.RedisConfig{
		Host:     "redis",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	redisDb, err := store.NewRedis(redisCfg)
	if err != nil {
		// log.Error().Err().Msg()
		log.Fatalf("failed to open redis conn: %v\n", err)
	}

	serverDataSource := &api.DataSource{
		PostgresDB: postgresDB,
		RedisDB:    redisDb,
	}

	srv := api.NewServer(serverCfg, serverDataSource)
	srv.Start()
}
