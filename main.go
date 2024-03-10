package main

import (
	"context"
	"fmt"
	"trade-balance-service/app"
	"trade-balance-service/config"
	// "github.com/caarlos0/env"
)

func main() {
	var cfg config.Config
	// if err := env.Parse(&cfg); err != nil {
	// 	return
	// }

	cfg.RabbitUser = "admin"
	cfg.RabbitPassword = "admin"
	cfg.RabbitHost = "localhost"
	cfg.RabbitPort = "5672"

	cfg.PostgreeUser = "admin"
	cfg.PostgreePassword = "admin"
	cfg.PostgreeHost = "localhost"
	cfg.PostgreePort = "5432"
	cfg.PostgreeDB = "bps"

	app.StartProgram(context.Background(), buildPostgreeUrl(cfg), buildRabbitUrl(cfg))
}

func buildRabbitUrl(cfg config.Config) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.RabbitUser, cfg.RabbitPassword, cfg.RabbitHost, cfg.RabbitPort)
}

func buildPostgreeUrl(cfg config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgreeUser, cfg.PostgreePassword, cfg.PostgreeHost, cfg.PostgreePort, cfg.PostgreeDB)
}
