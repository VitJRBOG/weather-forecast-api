package server

import (
	"database/sql"
	"log"
	"net/http"
)

func handle(dbConn *sql.DB) {
	http.HandleFunc("/cities", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			_, err := w.Write([]byte("List of cities"))
			if err != nil {
				log.Println(err)
			}
			// TODO
		default:
			_, err := w.Write([]byte("Method not allowed"))
			if err != nil {
				log.Println(err)
			}
			// TODO
		}
	})

	http.HandleFunc("/forecasts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			_, err := w.Write([]byte("List of forecasts for selected city"))
			if err != nil {
				log.Println(err)
			}
			// TODO
		default:
			_, err := w.Write([]byte("Method not allowed"))
			if err != nil {
				log.Println(err)
			}
			// TODO
		}
	})

	http.HandleFunc("/forecast", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			_, err := w.Write([]byte("Full forecast info by selected date"))
			if err != nil {
				log.Println(err)
			}
			// TODO
		default:
			_, err := w.Write([]byte("Method not allowed"))
			if err != nil {
				log.Println(err)
			}
			// TODO
		}
	})
}
