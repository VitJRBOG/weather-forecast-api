package db

import (
	"database/sql"
	"log"
	"time"
	"weather-forecast-api/internal/models"

	_ "github.com/lib/pq" // Postgres driver
)

func NewConnection(dsn string) *sql.DB {
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("unable connect to database: %s\n", err.Error())
	}

	return dbConn
}

func SelectAllFromCities(dbConn *sql.DB) []models.City {
	query := "SELECT * FROM cities"

	rows, err := dbConn.Query(query)
	if err != nil {
		log.Println(err)
		return nil
	}

	cities := []models.City{}

	for rows.Next() {
		city := models.City{}

		if err := rows.Scan(&city.ID, &city.Name, &city.Country,
			&city.Latitude, &city.Longitude); err != nil {
			log.Println(err)
			return nil
		}

		cities = append(cities, city)
	}

	return cities
}

func SelectByCityAndDateFromForecast(dbConn *sql.DB, cityID int, date time.Time) []models.Forecast {
	query := "SELECT * FROM forecast WHERE city_id = $1 AND f_date = $2"

	rows, err := dbConn.Query(query, cityID, date)
	if err != nil {
		log.Println(err)
		return nil
	}

	forecasts := []models.Forecast{}

	for rows.Next() {
		forecast := models.Forecast{}
		if err := rows.Scan(&forecast.ID, &forecast.Temperature, &forecast.Date,
			&forecast.FullInfo, &forecast.CityID); err != nil {
			log.Println(err)
			return nil
		}

		forecasts = append(forecasts, forecast)
	}

	return forecasts
}

func InsertIntoForecast(dbConn *sql.DB, forecast models.Forecast) {
	query := "INSERT INTO forecast(temp, f_date, full_info, city_id) VALUES($1, $2, $3, $4)"

	_, err := dbConn.Exec(query, forecast.Temperature, forecast.Date,
		forecast.FullInfo, forecast.CityID)
	if err != nil {
		log.Println(err)
		return
	}
}

func UpdateRawInForecast(dbConn *sql.DB, forecastID int, forecast models.Forecast) {
	query := "UPDATE forecast SET temp = $1, full_info = $2 WHERE id = $3"

	_, err := dbConn.Exec(query, forecast.Temperature, forecast.FullInfo, forecastID)
	if err != nil {
		log.Println(err)
		return
	}
}
