package models

import (
	"time"
)

type Forecast struct {
	ID          int
	Temperature float64
	Date        time.Time
	FullInfo    []byte
	CityID      int
}
