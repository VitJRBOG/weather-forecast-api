package weatherforecastapi

import (
	"fmt"
	"weather-forecast-api/internal/config"
	"weather-forecast-api/internal/db"
)

func Execute() {
	dbConnectionCfg := config.NewDBConnectionCfg()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConnectionCfg.User, dbConnectionCfg.Password,
		dbConnectionCfg.HostAddress, dbConnectionCfg.HostPort,
		dbConnectionCfg.DBName,
		dbConnectionCfg.SSLMode)

	db.NewConnection(dsn)

	// ...
}
