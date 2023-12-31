package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/metafates/smartway-test/config"
	v1 "github.com/metafates/smartway-test/internal/controller/http/v1"
	"github.com/metafates/smartway-test/internal/usecase"
	"github.com/metafates/smartway-test/internal/usecase/repository"
	"github.com/metafates/smartway-test/pkg/httpserver"
	"github.com/metafates/smartway-test/pkg/logger"
	"github.com/metafates/smartway-test/pkg/postgres"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	pg, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		l.Fatal(err)
	}
	defer pg.Close()

	repo := repository.NewPostgresRepository(pg)
	useCases := usecase.UseCases{
		Account:  usecase.NewAccountUseCase(repo),
		Schema:   usecase.NewSchemaUseCase(repo),
		Provider: usecase.NewProviderUseCase(repo),
		Airline:  usecase.NewAirlineUseCase(repo),
	}

	// HTTP server
	router := mux.NewRouter()
	router.Use(handlers.RecoveryHandler())
	router.Use(func(handler http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stderr, handler)
	})

	v1.RegisterRoutes(router, useCases, l)
	httpServer := httpserver.New(
		router,
		httpserver.Port(cfg.HTTP.Port),
	)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
