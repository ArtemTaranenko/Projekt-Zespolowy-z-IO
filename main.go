package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"english-learning-app/db"
)

func main() {
	seed := flag.Bool("seed", false, "Dodaj przykładowe quizy do pustej tabeli")
	flag.Parse()
	database, err := db.Open("./db/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	if *seed {
		if err := db.Seed(database); err != nil {
			log.Fatal(err)
		}
	}
	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           newRouter(database),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Println("API: http://127.0.0.1:8080/api/health")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
