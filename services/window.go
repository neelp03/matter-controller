package services

import (
    "fmt"
    "net/http"
    "os"
    "time"
		"sync"

    "github.com/neelp03/matter-controller/handlers"
    "github.com/neelp03/matter-controller/services"
    "github.com/neelp03/matter-controller/utils"
)

var WindowOpen = false
var WindowEventFlag = false
var WindowMu sync.Mutex

func GetWindowStatus() (bool, bool) {
	WindowMu.Lock()
	defer WindowMu.Unlock()
	if WindowEventFlag {
		WindowEventFlag = false
		return WindowOpen, true
	}
	return WindowOpen, false
}

func UpdateWindowStatus() {
	WindowMu.Lock()
	WindowOpen = !WindowOpen
	WindowEventFlag = !WindowEventFlag
	status := "closed"
	if WindowOpen {
		status = "open"
	}
	WindowMu.Unlock()
}
// evaluateRules is the rule‑based controller responsible for deciding when
// to open or close the window.  
//
// Inputs (read every 30 s)
//   • Indoor temperature (°C)
//   • Outdoor temperature (°C)
//   • Rain rate (mm)              – from Open‑Meteo Weather API
//   • Outdoor air‑quality index   – US AQI consolidated value from
//                                    Open‑Meteo Air‑Quality API
//
// Decision logic
// 1. Indoor > maxComfort and
//    (indoor − outdoor) ≥ deltaThreshold.
// 2. Safe conditions: Rain ≤ rainThreshold and AQI ≤ aqiThreshold.
// 3. If both 1 & 2 hold and the window is closed ⇒ OPEN.
// 4. If the window is open and any of the following are true ⇒ CLOSE:
//    • Indoor < minComfort.
//    • Outdoor ≥ indoor (no cooling advantage).
//    • Rain > rainThreshold.
//    • AQI > aqiThreshold.	

func evaluateRules() {
	const (
			minComfort     = 68.0
			maxComfort     = 77.0 
			deltaThreshold = 3.6
			rainThreshold  = 0.0
			aqiThreshold   = 100.0
			tick           = 30 * time.Second
	)

	for {
			indoor, errIn := services.ReadTemperature()

			outdoor, errOut := services.FetchOutdoorTemperature()
			rain, errRain := services.FetchOutdoorRain()
			aqi, errAQI := services.FetchOutdoorAQI()

			if errIn != nil || errOut != nil || errRain != nil || errAQI != nil {
					fmt.Println("[rules] sensor/API error; ensuring window closed and retrying:", errIn, errOut, errRain, errAQI)
					services.WindowMu.Lock()
					open := services.WindowOpen
					services.WindowMu.Unlock()
					if open {
							_ = services.CloseWindow()
					}
					time.Sleep(tick)
					continue
			}

			services.WindowMu.Lock()
			open := services.WindowOpen
			services.WindowMu.Unlock()

			ventilationAdvantage := indoor > maxComfort && (indoor-outdoor) >= deltaThreshold
			safeConditions := rain <= rainThreshold && aqi <= aqiThreshold

			switch {
			case ventilationAdvantage && safeConditions && !open:
					fmt.Printf("[rules] +++++ Opening window (indoor %.1f °F, outdoor %.1f °F, rain %.2f mm, AQI %.0f)\n", indoor, outdoor, rain, aqi)
					if err := services.OpenWindow(); err != nil {
							fmt.Println("[rules] failed to open window:", err)
					}

			case open && (!safeConditions || indoor < minComfort || outdoor >= indoor):
					fmt.Printf("[rules] ----- Closing window (indoor %.1f °F, outdoor %.1f °F, rain %.2f mm, AQI %.0f)\n", indoor, outdoor, rain, aqi)
					if err := services.CloseWindow(); err != nil {
							fmt.Println("[rules] failed to close window:", err)
					}
			}

			time.Sleep(tick)
	}
}

func OpenWindow() error {
	// Simulate opening the window
	fmt.Println("Opening window...")
	// UpdateWindowStatus()
	return nil
}

func CloseWindow() error {
	// Simulate closing the window
	fmt.Println("Closing window...")
	// UpdateWindowStatus()
	return nil
}