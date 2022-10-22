package models

import (
	"time"
)

type Forecast struct {
	ID          int       `json:"id"`
	Temperature float64   `json:"temperature"`
	Date        time.Time `json:"date"`
	FullInfo    []byte    `json:"full_info"`
	CityID      int       `json:"city_id"`
}
