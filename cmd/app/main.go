package main

import (
	"log"
	"os"

	"github.com/metafates/smartway-test/config"
	"github.com/metafates/smartway-test/internal/app"
)

func main() {
	os.Setenv("PG_URL", "1")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
