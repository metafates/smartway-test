package main

import (
	"log"

	"github.com/metafates/smartway-test/config"
	"github.com/metafates/smartway-test/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
