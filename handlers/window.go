package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var windowOpen = false
var windowMu sync.Mutex

func GetWindowStatus() string {
	windowMu.Lock()
	defer windowMu.Unlock()
	if windowOpen {
		return "open"
	}
	return "closed"
}

func WindowStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := GetWindowStatus()
	response := map[string]string{"window": status}
	json.NewEncoder(w).Encode(response)
}

func ToggleWindowHandler(w http.ResponseWriter, r *http.Request) {
	windowMu.Lock()
	windowOpen = !windowOpen
	status := "closed"
	if windowOpen {
		status = "open"
	}
	windowMu.Unlock()
	fmt.Println("========== Toggled window. Now:", status, "==========")
	response := map[string]string{"window": status}
	json.NewEncoder(w).Encode(response)
}
