package feeds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WeatherData represents weather information
type WeatherData struct {
	Summary      string  `json:"summary"`
	TemperatureC float64 `json:"temperatureC"`
	FeelsLikeC   float64 `json:"feelsLikeC"`
}

// Coordinates represents latitude and longitude
type Coordinates struct {
	Lat float64
	Lon float64
}

// OpenMeteoResponse represents the API response from Open-Meteo
type OpenMeteoResponse struct {
	Current struct {
		Temperature         float64 `json:"temperature_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		WeatherCode         int     `json:"weather_code"`
	} `json:"current"`
}

// Asia country coordinates (major cities)
var asiaCountryCoordinates = map[string]Coordinates{
	"JP": {Lat: 35.6762, Lon: 139.6503}, // Tokyo
	"CN": {Lat: 39.9042, Lon: 116.4074}, // Beijing
	"IN": {Lat: 28.6139, Lon: 77.2090},  // New Delhi
	"SG": {Lat: 1.3521, Lon: 103.8198},  // Singapore
	"HK": {Lat: 22.3193, Lon: 114.1694}, // Hong Kong
	"KR": {Lat: 37.5665, Lon: 126.9780}, // Seoul
	"TH": {Lat: 13.7563, Lon: 100.5018}, // Bangkok
	"ID": {Lat: -6.2088, Lon: 106.8456}, // Jakarta
	"MY": {Lat: 3.1390, Lon: 101.6869},  // Kuala Lumpur
	"PH": {Lat: 14.5995, Lon: 120.9842}, // Manila
	"VN": {Lat: 21.0285, Lon: 105.8542}, // Hanoi
	"TW": {Lat: 25.0330, Lon: 121.5654}, // Taipei
}

// Weather code to description mapping (WMO Weather interpretation codes)
var weatherCodeDescriptions = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Foggy",
	48: "Depositing rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	77: "Snow grains",
	80: "Slight rain showers",
	81: "Moderate rain showers",
	82: "Violent rain showers",
	85: "Slight snow showers",
	86: "Heavy snow showers",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

// FetchWeather fetches weather data for a given country using Open-Meteo API
func FetchWeather(country string) (*WeatherData, error) {
	// Get coordinates for the country
	coords, ok := asiaCountryCoordinates[country]
	if !ok {
		// Default to Tokyo if country not found
		coords = asiaCountryCoordinates["JP"]
	}

	// Build Open-Meteo API URL
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m,apparent_temperature,weather_code",
		coords.Lat, coords.Lon,
	)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make API request
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("weather API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	// Parse response
	var apiResp OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse weather response: %w", err)
	}

	// Convert weather code to description
	description, ok := weatherCodeDescriptions[apiResp.Current.WeatherCode]
	if !ok {
		description = "Unknown"
	}

	return &WeatherData{
		Summary:      description,
		TemperatureC: apiResp.Current.Temperature,
		FeelsLikeC:   apiResp.Current.ApparentTemperature,
	}, nil
}
