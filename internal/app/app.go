package app

import (
	"github.com/metafates/smartway-test/config"
	"github.com/metafates/smartway-test/pkg/logger"
)

type App struct {
	logger *logger.Logger
}

func New(cfg *config.Config) (*App, error) {
	app := &App{
		logger: logger.New(cfg.Log.Level),
	}

	return app, nil
}

func (a *App) Run() {
}
