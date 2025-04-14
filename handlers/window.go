package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/neelp03/matter-controller/services"
)

func WindowStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := GetWindowStatus()
	response := map[string]string{"window": status}
	json.NewEncoder(w).Encode(response)
}

func ToggleWindowHandler(w http.ResponseWriter, r *http.Request) {
	services.WindowMu.Lock()
	services.WindowOpen = !windowOpen
	services.WindowEventFlag = true
	status := "closed"
	if windowOpen {
		status = "open"
	}
	windowMu.Unlock()
	fmt.Println("========== Toggled window. Now:", status, "==========")
	response := map[string]string{"window": status}
	json.NewEncoder(w).Encode(response)
}
