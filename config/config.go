package config

import (
	_ "embed"
	"errors"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

//go:embed config.toml
var configFile []byte

type Config struct {
	HTTP     HTTPConfig `koanf:"http"`
	Log      LogConfig  `koanf:"log"`
	Postgres Postgres   `koanf:"postgres"`
}

type Postgres struct {
	URL     string `koanf:"url" validate:"required"`
	PoolMax int    `koanf:"poolmax"`
}

type HTTPConfig struct {
	Port string `koanf:"port"`
}

type LogConfig struct {
	Level string `koanf:"level"`
}

func New() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		if err = godotenv.Load(".env"); err != nil {
			return nil, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	k := koanf.New(".")

	err := k.Load(rawbytes.Provider(configFile), toml.Parser())
	if err != nil {
		return nil, err
	}

	const envPrefix = "SERVER_"
	err = k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.ReplaceAll(
			strings.ToLower(strings.TrimPrefix(s, envPrefix)),
			"_",
			".",
		)
	}), nil)

	if err != nil {
		return nil, err
	}

	var config Config
	err = k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{
		Tag:       "koanf",
		FlatPaths: false,
	})
	if err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return &config, nil
}
