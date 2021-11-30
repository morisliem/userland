package main

import (
	"context"
	"os"
	"strconv"
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

	server_readTimeout, _ := strconv.Atoi(os.Getenv("READ_TIMEOUT"))
	server_writeTimeout, _ := strconv.Atoi(os.Getenv("WRITE_TIMEOUT"))
	server_shutdownTimeout, _ := strconv.Atoi(os.Getenv("SHUTDOWN_TIMEOUT"))
	serverCfg := api.ServerConfig{
		Host:            os.Getenv("SERVER_HOST"),
		Port:            os.Getenv("SERVER_PORT"),
		ReadTimeout:     time.Duration(server_readTimeout) * time.Millisecond,
		WriteTimeout:    time.Duration(server_writeTimeout) * time.Millisecond,
		ShutdownTimeout: time.Duration(server_shutdownTimeout) * time.Second,
	}

	db_port, _ := strconv.Atoi(os.Getenv("PGPORT"))
	postgresCfg := store.PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     db_port,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}
	postgresDB, err := store.NewPG(postgresCfg)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}

	logStore := postgres.NewLoginStore(postgresDB)

	redis_port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	redis_db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisCfg := store.RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     redis_port,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redis_db,
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
		"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVER"),
	}

	consCfg := &kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVER"),
		"group.id":          os.Getenv("GROUP_ID"),
		"auto.offset.reset": os.Getenv("AUTO_OFFSET"),
	}

	msgBroker, err := broker.NewBroker(consCfg, prodCfg)
	if err != nil {
		panic(err)
	}

	msgBroker.CreateTopic(context.Background())

	terminateWorkerChan := make(chan int, 1)
	defer func() {
		terminateWorkerChan <- 1
	}()

	go kafka_job.LoginLog(context.Background(), msgBroker, logStore, terminateWorkerChan)
	srv := api.NewServer(serverCfg, serverDataSource, msgBroker)
	srv.Start()
}
