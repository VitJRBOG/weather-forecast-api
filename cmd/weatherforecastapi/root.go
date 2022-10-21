package weatherforecastapi

import (
	"fmt"
	"log"
	"weather-forecast-api/internal/config"
	"weather-forecast-api/internal/db"
)

func Execute() {
	initializeLogger()

	dbConnectionCfg := config.NewDBConnectionCfg()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConnectionCfg.User, dbConnectionCfg.Password,
		dbConnectionCfg.HostAddress, dbConnectionCfg.HostPort,
		dbConnectionCfg.DBName,
		dbConnectionCfg.SSLMode)

	db.NewConnection(dsn)

	// ...
}

func initializeLogger() {
	log.SetFlags(log.Ldate | log.Llongfile)
}
