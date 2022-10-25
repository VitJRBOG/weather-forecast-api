package weatherforecastapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
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
	collectingData(dbConn, openWeatherAPIConnectionCfg.APIID)
}

func collectingData(dbConn *sql.DB, openWeatherAPIID string) {
	cities, err := db.SelectAllFromCities(dbConn)
	if err != nil {
		log.Fatalln("forecasts updating has been interrupted")
	}

	wg := &sync.WaitGroup{}
	for _, city := range cities {
		wg.Add(1)
		go fetchAndSave(wg, dbConn, openWeatherAPIID, city)
	}

	wg.Wait()
}

func fetchAndSave(wg *sync.WaitGroup, dbConn *sql.DB, openWeatherAPIID string, city models.City) {
	forecasts, err := openweather.FetchForecast(openWeatherAPIID,
		city.Latitude, city.Longitude)
	if err != nil {
		log.Fatalln("unable fetch data from OpenWeather")
	}

	saveDataToDB(dbConn, city, forecasts)

	wg.Done()
}

func saveDataToDB(dbConn *sql.DB, city models.City, forecasts []openweather.Forecast) {
	for _, forecast := range forecasts {
		data, err := json.Marshal(forecast.FullInfo)
		if err != nil {
			log.Fatalln(err)
		}

		f := models.Forecast{
			Temperature: forecast.Temperature,
			Date:        int64(forecast.Date),
			FullInfo:    data,
			CityID:      city.ID,
		}

		fs, err := db.SelectByCityAndDateFromForecast(dbConn, f.CityID, f.Date)
		if err != nil {
			log.Fatalln("unable save data to database")
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
