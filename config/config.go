package config

type Config struct {
	Rabbit struct {
		User     string
		Password string
		Host     string
		Port     string
	}
	Postgres struct {
		User     string
		Password string
		Host     string
		Port     string
		Db       string
	}
}
