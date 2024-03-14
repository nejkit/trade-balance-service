package config

import (
	"github.com/spf13/viper"
	"os"
)

func GetConfig() (*Config, error) {

	cManager := viper.New()

	cManager.AutomaticEnv()

	cManager.SetConfigFile("config.yml")

	if os.Getenv("DOCKER_COMPOSE") == "true" {
		cManager.SetConfigFile("configs/bps.yml")
	}

	if err := cManager.ReadInConfig(); err != nil {
		return nil, err
	}

	var conf Config

	if err := cManager.Unmarshal(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
