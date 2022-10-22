package weatherforecastapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"weather-forecast-api/internal/config"
	"weather-forecast-api/internal/data/openweather"
	"weather-forecast-api/internal/db"
	"weather-forecast-api/internal/models"
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

	TestFetchingData(dbConn, openWeatherAPIConnectionCfg)

	// ...
}

// FIXME: удалить после описания сервера
func TestFetchingData(dbConn *sql.DB, openWeatherAPIConnectionCfg config.OpenWeatherAPIConnectionCfg) {
	cities := db.SelectAllFromCities(dbConn)

	for _, city := range cities {
		forecasts := openweather.FetchForecast(openWeatherAPIConnectionCfg.APIID,
			city.Latitude, city.Longitude)

		TestSaveDataToDB(dbConn, city, forecasts)
	}
}

// FIXME: удалить после описания сервера
func TestSaveDataToDB(dbConn *sql.DB, city models.City, forecasts []openweather.Forecast) {
	for _, forecast := range forecasts {
		data, err := json.Marshal(forecast.FullInfo)
		if err != nil {
			panic(err)
		}

		f := models.Forecast{
			Temperature: forecast.Temperature,
			Date:        time.Unix(int64(forecast.Date), 0),
			FullInfo:    data,
			CityID:      city.ID,
		}

		fs := db.SelectByCityAndDateFromForecast(dbConn, f.CityID, f.Date)
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
