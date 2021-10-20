package main

import (
	"log"
	"time"
	"userland/api"
	"userland/store"
)

func main() {
	// TODO use external config management (toml?)
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

	serverDataSource := &api.DataSource{
		PostgresDB: postgresDB,
	}

	srv := api.NewServer(serverCfg, serverDataSource)
	srv.Start()
}
