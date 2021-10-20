package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userland/api/handler/auth"
	"userland/store"
	"userland/store/postgres"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Config     ServerConfig
	DataSource *DataSource
	stores     *stores
}

type DataSource struct {
	PostgresDB *sql.DB
}

type stores struct {
	userStore store.UserStore
}

type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// init postgresql with lib (with standard, no handling connection directly?)
func NewServer(config ServerConfig, dataSource *DataSource) *Server {
	return &Server{
		Config:     config,
		DataSource: dataSource,
	}
}

func (s *Server) Start() {
	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(osSigChan)
		os.Exit(0)
	}()

	_ = s.initStores()

	r := s.createHandlers()
	address := fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port)
	srv := &http.Server{
		Addr:         address,
		ReadTimeout:  s.Config.ReadTimeout,
		WriteTimeout: s.Config.WriteTimeout,
		Handler:      r,
	}

	shutdownCtx := context.Background()
	if s.Config.ShutdownTimeout > 0 {
		var cancelShutdownTimeout context.CancelFunc
		shutdownCtx, cancelShutdownTimeout = context.WithTimeout(shutdownCtx, s.Config.ShutdownTimeout)
		defer cancelShutdownTimeout()
	}

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		// TODO replace with zlogger fatal
		panic("cannot start server")
	}
	// TODO with proper logging with zlogger
	fmt.Printf("serving %s\n", address)

	go func(srv *http.Server) {
		<-osSigChan
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			panic("failed to shutdown gracefully")
		}
	}(srv)
}

func (s *Server) initStores() error {
	userStore := postgres.NewUserStore(s.DataSource.PostgresDB)
	s.stores = &stores{
		userStore: userStore,
	}
	return nil
}

func (s *Server) createHandlers() http.Handler {
	// TODO pprof and healthcheck
	r := chi.NewRouter()
	r.Post("/register", auth.Register(s.stores.userStore))
	r.Post("/login", auth.Login(s.stores.userStore))

	return r
}
