package weatherforecastapi

import (
	"fmt"
	"log"
	"weather-forecast-api/internal/config"
	"weather-forecast-api/internal/data/openweather"
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

	openWeatherAPIConnectionCfg := config.NewOpenWeatherAPIConnectionCfg()

	TestFetchingData(openWeatherAPIConnectionCfg)

	// ...
}

// FIXME: удалить после описания работы с БД
func TestFetchingData(openWeatherAPIConnectionCfg config.OpenWeatherAPIConnectionCfg) {
	lat, lon := 55.750446, 37.617494

	fmt.Println(openweather.FetchForecast(openWeatherAPIConnectionCfg.APIID, lat, lon))
}

func initializeLogger() {
	log.SetFlags(log.Ldate | log.Llongfile)
}
