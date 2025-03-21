package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
)

type AddGeoResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func GetCoordinatesFromIP(ip string) (AddGeoResponse, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return AddGeoResponse{}, err
	}
	defer resp.Body.Close()

	var data AddGeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return AddGeoResponse{}, err
	}
	return data, nil
}

type GeoResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func GetCoordinatesFromCity(cityName string) (AddGeoResponse, error) {
	if cityName == "" {
		return AddGeoResponse{}, fmt.Errorf("city is empty %s", cityName)
	}
	baseURL := "https://nominatim.openstreetmap.org/search"
	query := url.QueryEscape(cityName)
	fullURL := fmt.Sprintf("%s?q=%s&format=json&limit=1", baseURL, query)

	resp, err := http.Get(fullURL)
	if err != nil {
		return AddGeoResponse{}, err
	}
	defer resp.Body.Close()

	var results []GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return AddGeoResponse{}, err
	}

	if len(results) == 0 {
		return AddGeoResponse{}, fmt.Errorf("can't found %s", cityName)
	}

	lat, lon := parseCoordinate(results[0].Lat), parseCoordinate(results[0].Lon)
	return AddGeoResponse{Lat: lat, Lon: lon}, nil
}

func parseCoordinate(coord string) float64 {
	var value float64
	fmt.Sscanf(coord, "%f", &value)
	return value
}

func CalculateDistance(add1, add2 AddGeoResponse) float64 {
	const earthRadius = 6371
	dLat := degreesToRadians(add2.Lat - add1.Lat)
	dLon := degreesToRadians(add2.Lon - add2.Lon)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(add1.Lat))*math.Cos(degreesToRadians(add2.Lat))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

type OSMGeoResponse struct {
	DisplayName string `json:"display_name"`
}

func GetAddressFromLatLonOSM(add AddGeoResponse) (string, error) {
	baseURL := "https://nominatim.openstreetmap.org/reverse"
	queryParams := url.Values{}
	queryParams.Set("lat", fmt.Sprintf("%f", add.Lat))
	queryParams.Set("lon", fmt.Sprintf("%f", add.Lon))
	queryParams.Set("format", "json")

	fullURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var geoResponse OSMGeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResponse); err != nil {
		return "", err
	}

	if geoResponse.DisplayName == "" {
		return "", fmt.Errorf("No address found for coordinates (%f, %f)", add.Lat, add.Lon)
	}

	return geoResponse.DisplayName, nil
}
