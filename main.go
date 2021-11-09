package main

import (
	"context"
	"time"
	"userland/api"
	"userland/api/kafka_job"
	"userland/store"
	"userland/store/broker"
	"userland/store/postgres"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	logStore := postgres.NewLoginStore(postgresDB)

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

	prodCfg := &kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	}

	consCfg := &kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "userland",
		"auto.offset.reset": "latest",
	}

	msgBroker, err := broker.NewBroker(consCfg, prodCfg)
	if err != nil {
		panic(err)
	}

	terminateWorkerChan := make(chan int, 1)
	defer func() {
		terminateWorkerChan <- 1
	}()

	go kafka_job.LoginLog(context.Background(), msgBroker, logStore, terminateWorkerChan)

	srv := api.NewServer(serverCfg, serverDataSource, msgBroker)
	srv.Start()
}
