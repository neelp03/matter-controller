package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var windowOpen = false
var windowEventFlag = false
var windowMu sync.Mutex

func GetWindowStatus() (boolean, boolean) {
	windowMu.Lock()
	defer windowMu.Unlock()
	if windowEventFlag {
		windowEventFlag = false
		return windowOpen, true
	}
	return windowOpen, false
}

func WindowStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := GetWindowStatus()
	response := map[string]string{"window": status}
	json.NewEncoder(w).Encode(response)
}

func ToggleWindowHandler(w http.ResponseWriter, r *http.Request) {
	windowMu.Lock()
	windowOpen = !windowOpen
	windowEventFlag = true
	status := "closed"
	if windowOpen {
		status = "open"
	}
	windowMu.Unlock()
	fmt.Println("========== Toggled window. Now:", status, "==========")
	response := map[string]string{"window": status}
	json.NewEncoder(w).Encode(response)
}
