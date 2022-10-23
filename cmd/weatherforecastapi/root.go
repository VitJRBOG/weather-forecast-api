package weatherforecastapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"weather-forecast-api/internal/config"
	"weather-forecast-api/internal/data/openweather"
	"weather-forecast-api/internal/db"
	"weather-forecast-api/internal/models"
	"weather-forecast-api/internal/server"
)

func Execute() {
	initializeLogger()

	dbConnectionCfg := config.NewDBConnectionCfg()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConnectionCfg.User, dbConnectionCfg.Password,
		dbConnectionCfg.HostAddress, dbConnectionCfg.HostPort,
		dbConnectionCfg.DBName,
		dbConnectionCfg.SSLMode)

	dbConn := db.NewConnection(dsn)

	openWeatherAPIConnectionCfg := config.NewOpenWeatherAPIConnectionCfg()

	updateForecastsInfoInDB(dbConn, openWeatherAPIConnectionCfg)

	serverCfg := config.NewServerCfg()

	server.Up(dbConn, serverCfg)
}

func updateForecastsInfoInDB(dbConn *sql.DB, openWeatherAPIConnectionCfg config.OpenWeatherAPIConnectionCfg) {
	fetchingData(dbConn, openWeatherAPIConnectionCfg.APIID)
}

func fetchingData(dbConn *sql.DB, openWeatherAPIID string) {
	cities, err := db.SelectAllFromCities(dbConn)
	if err != nil {
		log.Fatalln("forecasts updating has been interrupted")
	}

	for _, city := range cities {
		forecasts := openweather.FetchForecast(openWeatherAPIID,
			city.Latitude, city.Longitude)

		saveDataToDB(dbConn, city, forecasts)
	}
}

func saveDataToDB(dbConn *sql.DB, city models.City, forecasts []openweather.Forecast) {
	for _, forecast := range forecasts {
		data, err := json.Marshal(forecast.FullInfo)
		if err != nil {
			log.Fatalln("forecasts updating has been interrupted")
		}

		f := models.Forecast{
			Temperature: forecast.Temperature,
			Date:        int64(forecast.Date),
			FullInfo:    data,
			CityID:      city.ID,
		}

		fs, err := db.SelectByCityAndDateFromForecast(dbConn, f.CityID, f.Date)
		if err != nil {
			log.Fatalln("forecasts updating has been interrupted")
		}

		if len(fs) == 0 {
			db.InsertIntoForecast(dbConn, f)
		} else {
			db.UpdateRawInForecast(dbConn, fs[0].ID, f)
		}
	}
}

func initializeLogger() {
	log.SetFlags(log.Ldate | log.Llongfile)
}
