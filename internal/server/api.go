package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sort"
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

	sort.SliceStable(cities, func(i, j int) bool {
		return cities[i].Name < cities[j].Name
	})

	return cities, nil
}

type cityForecasts struct {
	Country       string  `json:"country"`
	CityName      string  `json:"city_name"`
	AvTemperature float64 `json:"av_temperature"`
	ForecastDates []int64 `json:"forecast_dates"`
}

func newCityForecasts(cities []models.City, forecasts []models.Forecast) cityForecasts {
	avTemperature := 0.0
	dates := []int64{}
	timeNow := time.Now().Unix()

	for _, forecast := range forecasts {
		avTemperature += forecast.Temperature
		if forecast.Date > timeNow {
			dates = append(dates, forecast.Date)
		}
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i] < dates[j]
	})

	avTemperature /= float64(len(forecasts))

	return cityForecasts{
		Country:       cities[0].Country,
		CityName:      cities[0].Name,
		AvTemperature: avTemperature,
		ForecastDates: dates,
	}
}

func getForecasts(dbConn *sql.DB, params url.Values) (cityForecasts, error) {
	if !params.Has("city_id") {
		return cityForecasts{}, Error{http.StatusBadRequest, "'city_id' param is empty"}
	}

	cityID, err := strconv.Atoi(params.Get("city_id"))
	if err != nil {
		return cityForecasts{}, Error{http.StatusBadRequest, "'city_id' must be integer"}
	}

	cities, err := db.SelectByIDFromCity(dbConn, cityID)
	if err != nil {
		return cityForecasts{}, Error{http.StatusServiceUnavailable, "couldn't get city info"}
	}

	if len(cities) == 0 {
		return cityForecasts{}, Error{http.StatusNotFound, "no cities found"}
	}

	forecasts, err := db.SelectByCityFromForecast(dbConn, cityID)
	if err != nil {
		return cityForecasts{}, Error{http.StatusServiceUnavailable, "couldn't get forecasts"}
	}

	return newCityForecasts(cities, forecasts), nil
}

type detailForecast struct {
	Temperature float64                  `json:"temperature"`
	FullInfo    []map[string]interface{} `json:"full_info"`
}

func newDetailForecast(temperature float64, fullInfo []byte) (detailForecast, error) {
	f := []map[string]interface{}{}

	err := json.Unmarshal(fullInfo, &f)
	if err != nil {
		log.Println(err)
		return detailForecast{}, Error{http.StatusServiceUnavailable, "service unavailable"}
	}

	return detailForecast{
		Temperature: temperature,
		FullInfo:    f,
	}, nil
}

func getForecast(dbConn *sql.DB, params url.Values) (detailForecast, error) {
	if !params.Has("city_id") {
		return detailForecast{}, Error{http.StatusBadRequest, "'city_id' param is empty"}
	}

	if !params.Has("date") {
		return detailForecast{}, Error{http.StatusBadRequest, "'date' param is empty"}
	}

	cityID, err := strconv.Atoi(params.Get("city_id"))
	if err != nil {
		return detailForecast{}, Error{http.StatusBadRequest, "'city_id' must be integer"}
	}

	date, err := strconv.ParseInt(params.Get("date"), 10, 64)
	if err != nil {
		return detailForecast{}, Error{http.StatusBadRequest, "'date' must be integer"}
	}

	forecasts, err := db.SelectByCityAndDateFromForecast(dbConn, cityID, date)
	if err != nil {
		return detailForecast{}, Error{http.StatusServiceUnavailable, "couldn't get forecast"}
	}

	if len(forecasts) == 0 {
		return detailForecast{}, Error{http.StatusNotFound, "no forecasts found"}
	}

	return newDetailForecast(forecasts[0].Temperature, forecasts[0].FullInfo)
}
