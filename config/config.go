package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

type Cfg struct {
	ServerPort int    `env:"ServerPort,default=8080"`
	DBHost     string `env:"DB_HOST,default=db"`
	DBPort     int    `env:"DB_PORT,default=5432"`
	DBUser     string `env:"DB_USER,default=postgres"`
	DBPass     string `env:"DB_PASS,default=topsecret"`
	DBName     string `env:"DB_NAME,default=pack_db"`
}

func New() *Cfg {
	var c Cfg
	if err := envdecode.Decode(&c); err != nil {
		log.Fatalf("Could not decode env variables: %v", err)
	}

	return &c
}
