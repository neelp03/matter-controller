package services

import (
	"fmt"
	"time"
)

/*
Rule based controller responsible for deciding if the window can be opened
based on the following inputs:
1. Rain rate (mm) from Open-Meteo Weather API
2. Outdoor air-quality index from Open-Meteo Air-Quality API
*/
func canWindowBeOpen() bool {}

func ruleBasedControllerEval() (bool, bool) {}

func modelBasedControllerEval() (bool, bool) {}

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
