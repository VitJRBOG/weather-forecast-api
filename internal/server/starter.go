package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"weather-forecast-api/internal/config"
)

func Up(dbConn *sql.DB, serverCfg config.ServerCfg) {
	handle(dbConn)

	address := fmt.Sprintf(":%s", serverCfg.Port)
	err := http.ListenAndServe(address, logging(http.DefaultServeMux))
	if err != nil {
		log.Fatalf("server error: %s\n", err.Error())
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begins := time.Now()
		next.ServeHTTP(w, r)
		timeElapsed := time.Since(begins)

		log.Printf("[%s] %s %s", r.Method, r.RequestURI, timeElapsed)
	})
}
