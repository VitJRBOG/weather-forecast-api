package models

type Forecast struct {
	ID          int     `json:"id"`
	Temperature float64 `json:"temperature"`
	Date        int64   `json:"date"`
	FullInfo    []byte  `json:"full_info"`
	CityID      int     `json:"city_id"`
}
