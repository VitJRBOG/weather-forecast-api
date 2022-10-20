package config

import (
	"log"
	"os"
)

type DBConnectionCfg struct {
	HostAddress string
	HostPort    string
	User        string
	Password    string
	DBName      string
	SSLMode     string
}

func NewDBConnectionCfg() DBConnectionCfg {
	cfg := DBConnectionCfg{}
	cfg.HostAddress = os.Getenv("POSTGRES_HOST_ADDRESS")
	cfg.HostPort = os.Getenv("POSTGRES_HOST_PORT")
	cfg.User = os.Getenv("POSTGRES_USER")
	cfg.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.DBName = os.Getenv("POSTGRES_DB")
	cfg.SSLMode = os.Getenv("SSL_MODE")

	someIsEmpty := false

	if cfg.HostAddress == "" {
		log.Println("POSTGRES_HOST_ADDRESS env variable is empty")
		someIsEmpty = true
	}

	if cfg.HostPort == "" {
		log.Println("POSTGRES_HOST_PORT env variable is empty")
		someIsEmpty = true
	}

	if cfg.User == "" {
		log.Println("POSTGRES_USER env variable is empty")
		someIsEmpty = true
	}

	if cfg.Password == "" {
		log.Println("POSTGRES_PASSWORD env variable is empty")
		someIsEmpty = true
	}

	if cfg.DBName == "" {
		log.Println("POSTGRES_DB env variable is empty")
		someIsEmpty = true
	}

	if cfg.SSLMode == "" {
		log.Println("SSL_MODE env variable is empty")
		someIsEmpty = true
	}

	if someIsEmpty {
		log.Fatalln("some desktop environments is empty")
	}

	return cfg
}
