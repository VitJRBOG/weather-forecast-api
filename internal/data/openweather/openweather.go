package openweather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Forecast struct {
	Date        int
	Temperature float64
	FullInfo    []fullForecastData
}

func FetchForecast(apiID string, latitude, longitude float64) []Forecast {
	u := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?units=metric&lat=%f&lon=%f&appid=%s",
		latitude, longitude, apiID)

	response, err := http.Get(u)
	if err != nil {
		log.Println(err)
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	forecast := parseAPIResponse(body)

	return forecast
}

type apiResponse struct {
	Code    string             `json:"cod"`
	Message any                `json:"message"`
	List    []fullForecastData `json:"list"`
}

type fullForecastData struct {
	Date       int              `json:"dt"`
	Main       mainForecastData `json:"main"`
	Weather    []weatherData    `json:"weather"`
	Clouds     cloudsData       `json:"clouds"`
	Wind       windData         `json:"wind"`
	Visibility int              `json:"visibility"`
	Pop        float64          `json:"pop"`
	Sys        sysData          `json:"sys"`
	DateTxt    string           `json:"dt_txt"`
}

type mainForecastData struct {
	Temperature float64 `json:"temp"`
	FeelsLike   float64 `json:"feels_like"`
	TempMin     float64 `json:"temp_min"`
	TempMax     float64 `json:"temp_max"`
	Pressure    int     `json:"pressure"`
	SeaLevel    int     `json:"sea_level"`
	GrndLevel   int     `json:"grnd_level"`
	Humidity    int     `json:"humidity"`
	TempKf      float64 `json:"temp_kf"`
}

type weatherData struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type cloudsData struct {
	All int `json:"all"`
}

type windData struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type sysData struct {
	Pod string `json:"pod"`
}

func parseAPIResponse(data []byte) []Forecast {
	apiResp := parseRawData(data)
	forecast := selectNecessaryInfo(apiResp.List)

	return forecast
}

func parseRawData(data []byte) apiResponse {
	apiResp := apiResponse{}

	err := json.Unmarshal(data, &apiResp)
	if err != nil {
		log.Println(err)
		return apiResponse{}
	}

	return apiResp
}

func selectNecessaryInfo(weatherData []fullForecastData) []Forecast {
	forecast := []Forecast{}

	beginDate := calculateBeginDate()

	for i := 1; i <= 5; i++ {
		endDate := beginDate + 86400

		wi := selectByDayForecast(beginDate, endDate, weatherData)

		forecast = append(forecast, wi)

		beginDate = endDate
	}

	return forecast
}

func calculateBeginDate() int {
	begin, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		log.Println(err)
		return 0
	}

	beginDate := begin.Unix()

	return int(beginDate)
}

func selectByDayForecast(beginDate, endDate int, forecastData []fullForecastData) Forecast {
	fullForecastDataByDay := []fullForecastData{}

	byHoursTemperature := []float64{}

	for _, item := range forecastData {
		if beginDate > item.Date || endDate <= item.Date {
			continue
		}

		byHoursTemperature = append(byHoursTemperature, item.Main.Temperature)
		fullForecastDataByDay = append(fullForecastDataByDay, item)
	}

	return Forecast{
		Date:        beginDate,
		Temperature: calculateAverageDayTemperature(byHoursTemperature),
		FullInfo:    fullForecastDataByDay,
	}
}

func calculateAverageDayTemperature(byHoursTemperature []float64) float64 {
	averageTemperature := 0.0

	for _, temperature := range byHoursTemperature {
		averageTemperature += temperature
	}

	return averageTemperature / float64(len(byHoursTemperature))
}
