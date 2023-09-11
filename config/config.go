package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Port string    `koanf:"port"`
	DB   DBConfig  `koanf:"db"`
	Log  LogConfig `koanf:"log"`
}

type DBConfig struct {
	PostgresURI string `koanf:"postgres"`
}

type LogConfig struct {
	Level string `koanf:"level"`
}

func Load(envFiles ...string) (config *Config, err error) {
	if err = godotenv.Load(envFiles...); err != nil {
		return
	}

	k := koanf.New(".")

	// Default values
	err = k.Load(confmap.Provider(map[string]any{
		"port": "1234",
	}, "."), nil)

	if err != nil {
		return
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
		return
	}

	err = k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{
		Tag:       "koanf",
		FlatPaths: false,
	})

	if err != nil {
		return
	}

	return
}
