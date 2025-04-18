package services

import (
	"fmt"
	"time"
)

const (
	deltaThreshold = 5.0
	targetComfort  = 70.0
	// These thresholds are arbitrary and should be adjusted based on user preferences.
	rainThreshold = 10.0
	aqiThreshold  = 100.0
	tick          = time.Minute
)

/*
Rule based check responsible for deciding if the window can be opened
based on the following inputs:
1. Rain rate (mm) from Open-Meteo Weather API
2. Outdoor air-quality index (US AQI) from Open-Meteo Air-Quality API
*/
func canWindowBeOpen() bool {
	rain, err := FetchOutdoorRain()
	if err != nil {
		fmt.Println("Error fetching outdoor rain:", err)
		return false
	}

	aqi, err := FetchOutdoorAQI()
	if err != nil {
		fmt.Println("Error fetching outdoor AQI:", err)
		return false
	}

	return rain <= rainThreshold && aqi <= aqiThreshold
}

/*
ruleBasedController is the rule‑based controller responsible for deciding when
to open or close the window.

Inputs (read every 1 minute):
  - Indoor temperature (°F) from the indoor sensor
  - Outdoor temperature (°F) from the Open-Meteo Weather API
  - Forecast in the future

Decision logic:

  - If the indoor temperature is outside the comfort range (65°F to 75°F) and the
    outdoor temperature is in the direction of the comfort range then open the window.
    ex. Indoor: 50°F, Outdoor: 100°F both are outside the comfort range but
    the outdoor temperature is in the direction of the comfort range and will bring
    the indoor range toward the comfort range.

    Returns:

  - triggerEvent: true if a window command should be sent, false if it should not

  - shouldOpen: true if the window should be opened, false if it should be closed

  - If triggerEvent is false, shouldOpen is ignored.
*/
func ruleBasedControllerEval() (bool, bool) {
	indoor, errIn := ReadTemperature()
	if errIn != nil {
		fmt.Println("Error reading indoor temperature:", errIn)
		return false, false
	}

	comfortNeeds := indoor < (targetComfort-deltaThreshold) || indoor > (targetComfort+deltaThreshold)
	if !comfortNeeds {
		return false, false
	}

	outdoor, errOut := FetchOutdoorTemperature()
	if errOut != nil {
		fmt.Println("Error fetching outdoor temperature:", errOut)
		return false, false
	}

	if (indoor < (targetComfort-deltaThreshold) && outdoor > indoor) || (indoor > (targetComfort+deltaThreshold) && outdoor < indoor) {
		return true, true
	} else {
		return true, false
	}
}

func modelBasedControllerEval() (bool, bool) { return false, false }

/*
Once a minute will check if the window can be open.
If it can't and the window is open, close it.
If it can the system will evaluate if a window event should be triggered.
*/
func automatedControllerEval() {
	// Check if the window can be open
	canBeOpen := canWindowBeOpen()

	if !canBeOpen {
		if WindowOpen {
			err := CloseWindow()
			if err != nil {
				fmt.Println("Error closing window:", err)
			} else {
				fmt.Println("Window closed")
			}
		}
		return
	}

	// Evaluate if a window event should be triggered
	// modelBasedControllerEval()
	trigger, shouldOpen := ruleBasedControllerEval()

	if !trigger {
		return
	}

	if shouldOpen {
		err := OpenWindow()
		if err != nil {
			fmt.Println("Error opening window:", err)
		} else {
			fmt.Println("Window opened")
		}
	} else {
		err := CloseWindow()
		if err != nil {
			fmt.Println("Error closing window:", err)
		} else {
			fmt.Println("Window closed")
		}
	}
}

func RunAutomatedController() {
	// Run the automated controller every minute
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			automatedControllerEval()
		}
	}
}
