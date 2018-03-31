package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/marionauta/geobikes"
)

var endpoint = geobikes.BikeServer{
	Contract: contract,
	Token:    token,
}

func main() {
	http.HandleFunc("/", stations)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func stations(w http.ResponseWriter, r *http.Request) {
	stations, err := endpoint.AllStations()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(stations.IntoGeoJSON())
}
