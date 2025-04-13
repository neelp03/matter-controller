package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WeatherResponse struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		Rain        float64 `json:"rain"`
	} `json:"current"`
}

func FetchOutdoorWeather() (*WeatherResponse, error) {
	url := "https://api.open-meteo.com/v1/forecast?latitude=37.7749&longitude=-122.4194&current=temperature_2m,rain"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("!!!!!!!!!! failed to call Open-Meteo API !!!!!!!!!!: %v", err)
	}
	defer resp.Body.Close()

	var weather WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, fmt.Errorf("!!!!!!!!!! failed to decode Open-Meteo response !!!!!!!!!!: %v", err)
	}

	return &weather, nil
}
