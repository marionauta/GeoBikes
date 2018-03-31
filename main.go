package geobikes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL    = "https://api.jcdecaux.com/vls/v1/stations"
	parameters = "?contract=%s&apiKey=%s"
)

// The BikeServer type stores information about the API server.
type BikeServer struct {
	Contract string
	Token    string
}

// query formats the query part of the url.
func (e BikeServer) query() string {
	return fmt.Sprintf(parameters, e.Contract, e.Token)
}

// AllStations returns all stations from the API.
func (e BikeServer) AllStations() (Stations, error) {
	url := baseURL + e.query()
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var stations []Station
	json.NewDecoder(res.Body).Decode(&stations)
	res.Body.Close()

	return stations, nil
}

// GetStation return the station with the given number from the API.
func (e BikeServer) GetStation(number int) (Station, error) {
	url := fmt.Sprintf("%s/%d%s", baseURL, number, e.query())
	res, err := http.Get(url)
	if err != nil {
		return Station{}, err
	}

	var station Station
	json.NewDecoder(res.Body).Decode(&station)
	res.Body.Close()

	return station, nil
}

// Geometry represents a GeoJSON coordinate.
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// GeoStationProperties stores the relevant information of a station.
type GeoStationProperties struct {
	Number          int     `json:"number"`
	Name            string  `json:"name"`
	AvailableBikes  int     `json:"available_bikes"`
	AvailableStands int     `json:"available_stands"`
	Occupation      float64 `json:"occupation"`
}

// GeoStation represents a GeoJSON station.
type GeoStation struct {
	Type       string               `json:"type"`
	Geometry   Geometry             `json:"geometry"`
	Properties GeoStationProperties `json:"properties"`
}

// GeoCollection represents a GeoJSON 'FeatureCollection' object.
type GeoCollection struct {
	Type     string       `json:"type"`
	Features []GeoStation `json:"features"`
}

// Position stores a coordinate
type Position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Station represents an original API station
type Station struct {
	Number          int      `json:"number"`
	Name            string   `json:"name"`
	Position        Position `json:"position"`
	AvailableBikes  int      `json:"available_bikes"`
	AvailableStands int      `json:"available_bike_stands"`
}

// Occupation returns the station's occupation percentage [0, 1].
// If there are no available bikes or stands, it returns 1.
func (s Station) Occupation() float64 {
	if s.AvailableBikes == 0 && s.AvailableStands == 0 {
		return 1
	}

	total := float64(s.AvailableBikes + s.AvailableStands)
	return float64(s.AvailableBikes) / total
}

// IntoGeoJSON transforms a normal Station into a GeoJSON one.
func (s Station) IntoGeoJSON() GeoStation {
	return GeoStation{
		Type: "Feature",
		Geometry: Geometry{
			Type:        "Point",
			Coordinates: []float64{s.Position.Lng, s.Position.Lat},
		},
		Properties: GeoStationProperties{
			Number:          s.Number,
			Name:            s.Name,
			AvailableBikes:  s.AvailableBikes,
			AvailableStands: s.AvailableStands,
			Occupation:      s.Occupation(),
		},
	}
}

// Stations is just a Station slice.
type Stations []Station

// IntoGeoJSON transforms a slice of normal stations into a GeoCollection with
// GeoStations in its features.
func (ss Stations) IntoGeoJSON() GeoCollection {
	coll := GeoCollection{
		Type:     "FeatureCollection",
		Features: []GeoStation{},
	}

	for _, s := range ss {
		coll.Features = append(coll.Features, s.IntoGeoJSON())
	}

	return coll
}
