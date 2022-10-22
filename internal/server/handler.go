package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Error struct {
	HTTPStatus int
	Detail     string
}

func (e Error) Error() string {
	return fmt.Sprintf("status %d: %s", e.HTTPStatus, e.Detail)
}

func handle(dbConn *sql.DB) {
	http.HandleFunc("/cities", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			cities, err := getCities(dbConn)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, cities)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
		}
	})

	http.HandleFunc("/forecasts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				sendError(w, err)
				return
			}

			forecasts, err := getForecasts(dbConn, r.Form)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, forecasts)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
		}
	})

	http.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			err := r.ParseForm()
			if err != nil {
				log.Println(err)
				sendError(w, err)
				return
			}

			forecast, err := getForecast(dbConn, r.Form)
			if err != nil {
				sendError(w, err)
				return
			}

			sendData(w, forecast)
		default:
			sendError(w, Error{http.StatusMethodNotAllowed, "method not allowed"})
		}
	})
}

func sendData(w http.ResponseWriter, values interface{}) {
	response := map[string]interface{}{
		"status":   http.StatusOK,
		"response": values,
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		sendError(w, err)
		return
	}
}

func sendError(w http.ResponseWriter, reqError error) {
	response := map[string]interface{}{
		"status": http.StatusInternalServerError,
		"error":  "internal server error",
	}

	if errInfo, ok := reqError.(Error); ok {
		response["status"] = errInfo.HTTPStatus
		response["error"] = errInfo.Detail
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
