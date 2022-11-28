package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RestPort string `envconfig:"REST_PORT" default:"8080"`
	GrpcPort string `envconfig:"GRPC_PORT" default:"9090"`

	DBHost     string `envconfig:"MYSQL_DB_HOST" default:"svc-mix-id-1"`
	DBPort     string `envconfig:"MYSQL_DB_PORT" default:"3306"`
	DBUsername string `envconfig:"MYSQL_DB_USERNAME" default:"user"`
	DBPassword string `envconfig:"MYSQL_DB_PASSWORD" default:"usertest"`
	DBName     string `envconfig:"MYSQL_DB_NAME" default:"user"`

	// JWT
	JWTKey string `envconfig:"JWT_KEY" default:"user"`
}

func New() Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	return c
}
