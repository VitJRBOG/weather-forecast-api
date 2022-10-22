package server

import (
	"database/sql"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"weather-forecast-api/internal/db"
	"weather-forecast-api/internal/models"
)

func getCities(dbConn *sql.DB) ([]models.City, error) {
	cities, err := db.SelectAllFromCities(dbConn)
	if err != nil {
		return nil, Error{http.StatusServiceUnavailable, "couldn't get cities list"}
	}

	return cities, nil
}

func getForecasts(dbConn *sql.DB, params url.Values) ([]models.Forecast, error) {
	if !params.Has("city_id") {
		return nil, Error{http.StatusBadRequest, "'city_id' param is empty"}
	}

	cityID, err := strconv.Atoi(params.Get("city_id"))
	if err != nil {
		return nil, Error{http.StatusBadRequest, "'city_id' must be integer"}
	}

	cities, err := db.SelectByIDFromCity(dbConn, cityID)
	if err != nil {
		return nil, Error{http.StatusServiceUnavailable, "couldn't get city info"}
	}

	if len(cities) == 0 {
		return nil, Error{http.StatusNotFound, "no cities found"}
	}

	forecasts, err := db.SelectByCityFromForecast(dbConn, cityID)
	if err != nil {
		return nil, Error{http.StatusServiceUnavailable, "couldn't get forecasts"}
	}

	// TODO: описать вывод данных на основе задания
	// (инфу о городе + прогнозы погоды, отсортированные по дате)

	return forecasts, nil
}

func getForecast(dbConn *sql.DB, params url.Values) (models.Forecast, error) {
	if !params.Has("city_id") {
		return models.Forecast{}, Error{http.StatusBadRequest, "'city_id' param is empty"}
	}

	if !params.Has("date") {
		return models.Forecast{}, Error{http.StatusBadRequest, "'date' param is empty"}
	}

	cityID, err := strconv.Atoi(params.Get("city_id"))
	if err != nil {
		return models.Forecast{}, Error{http.StatusBadRequest, "'city_id' must be integer"}
	}

	date, err := time.Parse("2006-01-02", params.Get("date"))
	if err != nil {
		return models.Forecast{}, Error{http.StatusBadRequest, "'date' has the invalid format"}
	}

	forecasts, err := db.SelectByCityAndDateFromForecast(dbConn, cityID, date)
	if err != nil {
		return models.Forecast{}, Error{http.StatusServiceUnavailable, "couldn't get forecast"}
	}

	if len(forecasts) == 0 {
		return models.Forecast{}, Error{http.StatusNotFound, "no forecasts found"}
	}

	return forecasts[0], nil
}
