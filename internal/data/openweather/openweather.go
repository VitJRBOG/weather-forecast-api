package openweather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Forecast struct {
	Date        int64
	Temperature float64
	FullInfo    []fullForecastData
}

func FetchForecast(apiID string, latitude, longitude float64) ([]Forecast, error) {
	u := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?units=metric&lat=%f&lon=%f&appid=%s",
		latitude, longitude, apiID)

	response, err := http.Get(u)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	forecast, err := parseAPIResponse(body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return forecast, nil
}

type apiResponse struct {
	Code    interface{}        `json:"cod"`
	Message interface{}        `json:"message"`
	List    []fullForecastData `json:"list"`
}

type fullForecastData struct {
	Date       int64            `json:"dt"`
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

func parseAPIResponse(data []byte) ([]Forecast, error) {
	apiResp := parseRawData(data)

	if statCode, ok := apiResp.Code.(float64); ok {
		if statCode != 200 {
			if msg, ok := apiResp.Message.(string); ok {
				return nil, errors.New(msg)
			}
			return nil, fmt.Errorf("%v", apiResp.Message)
		}
	}

	forecast := selectNecessaryInfo(apiResp.List)

	return forecast, nil
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
	forecasts := []Forecast{}

	beginDate := calculateBeginDate()
	endDate := beginDate + (86400 * 5)

	for _, item := range weatherData {
		if beginDate > item.Date || endDate <= item.Date {
			continue
		}

		forecast := Forecast{
			item.Date,
			item.Main.Temperature,
			weatherData,
		}

		forecasts = append(forecasts, forecast)
	}

	return forecasts
}

func calculateBeginDate() int64 {
	begin, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		log.Println(err)
		return 0
	}

	beginDate := begin.Unix()

	return beginDate
}
