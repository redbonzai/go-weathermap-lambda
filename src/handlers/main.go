package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda" //nolint:goimports
	"net/http"
)

type WeatherRequest struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type WeatherResponse struct {
	Condition   string `json:"condition"`
	Description string `json:"description"`
	Temperature string `json:"temperature"`
}

// This struct is directly used for decoding the OpenWeather API response
type WeatherAPIResponse struct {
	Main struct {
		Temp float64 `json:"temp"` // Temperature in Kelvin
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`        // General condition (Rain, Snow, etc)
		Description string `json:"description"` // More detailed condition description
	} `json:"weather"`
}

const openWeatherAPIKey = "de235288a5f2a9fd5d865f75d2bec841"

func Handler(ctx context.Context, request WeatherRequest) (*WeatherResponse, error) {

	resp, err := fetchWeatherData(request.Latitude, request.Longitude)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Weather in %s, %s: %s, %s, %s\n", request.Latitude, request.Longitude, resp.Condition, resp.Description, resp.Temperature)

	return resp, nil
}

func fetchWeatherData(lat, lon string) (*WeatherResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s4&lon=%s&appid=%s", lat, lon, openWeatherAPIKey)
	resp, err := http.Get(url)
	fmt.Printf("THe GET Response: %v\n", resp)
	if err != nil {
		return nil, fmt.Errorf("error fetching weather data: %v", err)
	}
	defer resp.Body.Close()

	var apiResponse WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding weather data: %v", err)
	}

	tempCelsius := apiResponse.Main.Temp - 273.15 // Convert from Kelvin to Celsius
	temperature := "moderate"
	if tempCelsius < 10 {
		temperature = "cold"
	} else if tempCelsius > 25 {
		temperature = "hot"
	}

	fmt.Printf("The weather %v\n\n", apiResponse.Weather)

	if len(apiResponse.Weather) == 0 {
		return nil, fmt.Errorf("weather data is missing")
	}

	response := WeatherResponse{
		Condition:   apiResponse.Weather[0].Main,
		Description: apiResponse.Weather[0].Description,
		Temperature: temperature,
	}
	fmt.Printf("Weather in %s, %s: %s, %s, %s\n", lat, lon, response.Condition, response.Description, response.Temperature)
	return &response, nil
}

func main() {
	lambda.Start(Handler)
}
