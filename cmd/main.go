package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	flag "github.com/spf13/pflag"

	"github.com/marionauta/geobikes"
)

var endpoint = geobikes.BikeServer{
	Contract: contract,
	Token:    token,
}

func main() {
	port := flag.IntP("port", "p", 80, "Port to use")
	flag.Parse()

	http.HandleFunc("/", stations)
	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(addr, nil))
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
