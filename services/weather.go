// weather.go
package services

import (
    "encoding/json"
    "fmt"
    "math"
    "net/http"
    "time"

    "github.com/neelp03/matter-controller/utils"
)

// ---------------------------------------------------------------------------
// Weather (temperature & rain)
// ---------------------------------------------------------------------------

// WeatherResponse mirrors the JSON returned by the standard forecast endpoint.
type WeatherResponse struct {
    Current struct {
        Temperature float64 `json:"temperature_2m"`
        Rain        float64 `json:"rain"`
    } `json:"current"`
}

func FetchOutdoorWeather() (*WeatherResponse, error) {
    url := "https://api.open-meteo.com/v1/forecast?latitude=37.3387&longitude=-121.8853&current=temperature_2m,rain"
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to call Open‑Meteo API: %w", err)
    }
    defer resp.Body.Close()

    var weather WeatherResponse
    if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
        return nil, fmt.Errorf("failed to decode Open‑Meteo response: %w", err)
    }
    return &weather, nil
}

func FetchOutdoorTemperature() (float64, error) {
    weather, err := FetchOutdoorWeather()
    if err != nil {
        return 0, err
    }
    return utils.CToF(weather.Current.Temperature), nil
}

func FetchOutdoorRain() (float64, error) {
    weather, err := FetchOutdoorWeather()
    if err != nil {
        return 0, err
    }
    return weather.Current.Rain, nil
}

// ---------------------------------------------------------------------------
// Air quality – converts hourly PM2.5 & PM10 to US‑EPA AQI
// ---------------------------------------------------------------------------

type aqiResponse struct {
    Hourly struct {
        Time []string  `json:"time"`
        PM10 []float64 `json:"pm10"`
        PM25 []float64 `json:"pm2_5"`
    } `json:"hourly"`
}

// FetchOutdoorAQI fetches the latest PM2.5 & PM10 readings for today and
// returns the higher of their calculated US AQI values.
func FetchOutdoorAQI() (float64, error) {
    url := "https://air-quality-api.open-meteo.com/v1/air-quality?latitude=37.3387&longitude=-121.8853&hourly=pm10,pm2_5&timezone=America%2FLos_Angeles&forecast_days=1"

    resp, err := http.Get(url)
    if err != nil {
        return 0, fmt.Errorf("failed to call Open‑Meteo AQI API: %w", err)
    }
    defer resp.Body.Close()

    var data aqiResponse
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return 0, fmt.Errorf("failed to decode AQI response: %w", err)
    }

    n := len(data.Hourly.Time)
    if n == 0 {
        return 0, fmt.Errorf("no AQI data returned")
    }

    // Use newest hour (last index)
    pm25 := data.Hourly.PM25[n-1]
    pm10 := data.Hourly.PM10[n-1]

    aqiPM25 := calcAQI(pm25, pm25Breakpoints)
    aqiPM10 := calcAQI(pm10, pm10Breakpoints)

    return math.Max(aqiPM25, aqiPM10), nil
}

// --- breakpoint tables -----------------------------------------------------

type breakpoint struct{
    Clow, Chigh float64 // concentration range (µg/m³)
    Ilow, Ihigh float64 // AQI range
}

var pm25Breakpoints = []breakpoint{
    {0.0, 12.0, 0, 50},
    {12.1, 35.4, 51, 100},
    {35.5, 55.4, 101, 150},
    {55.5, 150.4, 151, 200},
    {150.5, 250.4, 201, 300},
    {250.5, 350.4, 301, 400},
    {350.5, 500.4, 401, 500},
}

var pm10Breakpoints = []breakpoint{
    {0, 54, 0, 50},
    {55, 154, 51, 100},
    {155, 254, 101, 150},
    {255, 354, 151, 200},
    {355, 424, 201, 300},
    {425, 504, 301, 400},
    {505, 604, 401, 500},
}

// calcAQI linearly interpolates the AQI for a pollutant concentration given
// its breakpoint table.
func calcAQI(c float64, table []breakpoint) float64 {
    for _, bp := range table {
        if c >= bp.Clow && c <= bp.Chigh {
            return (bp.Ihigh-bp.Ilow)/(bp.Chigh-bp.Clow)*(c-bp.Clow) + bp.Ilow
        }
    }
    return 500 // cap at hazardous if above table
}

// ---------------------------------------------------------------------------
// Utility – round timestamped sample arrays to latest hour (optional)
// ---------------------------------------------------------------------------

// latestIndex returns the index of the newest ISO‑8601 timestamp.
func latestIndex(times []string) int {
    latest := 0
    var latestTime time.Time
    for i, t := range times {
        if parsed, err := time.Parse(time.RFC3339, t+":00"); err == nil {
            if parsed.After(latestTime) {
                latestTime = parsed
                latest = i
            }
        }
    }
    return latest
}
