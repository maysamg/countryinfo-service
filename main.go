package main

import (
	"countryinfo-service/internal"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Sett oppstartstidspunkt
	internal.StartTime = time.Now() //  Lagre tidspunktet serveren startet

	// Sett opp port
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT ikke satt. Bruker standardport 8080.")
		port = "8080"
	}

	// Registrer endepunkter
	http.HandleFunc("/countryinfo/v1/info/", internal.InfoHandler)
	http.HandleFunc("/countryinfo/v1/population/", internal.PopulationHandler)
	http.HandleFunc("/countryinfo/v1/status/", internal.StatusHandler)

	// Start serveren
	log.Println("Starter serveren på port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
