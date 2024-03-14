package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"trade-balance-service/app"
	"trade-balance-service/config"
	// "github.com/caarlos0/env"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	cfg, err := config.GetConfig()

	if err != nil {
		logrus.Error(err.Error())
		return
	}

	app.StartProgram(context.Background(), buildPostgresUrl(cfg), buildRabbitUrl(cfg))
}

func buildRabbitUrl(cfg *config.Config) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.Rabbit.User, cfg.Rabbit.Password, cfg.Rabbit.Host, cfg.Rabbit.Port)
}

func buildPostgresUrl(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Db)
}
