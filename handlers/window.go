package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/neelp03/matter-controller/services"
)

func WindowStatusHandler(w http.ResponseWriter, r *http.Request) {
	status, _ := services.GetWindowStatus()
	response := map[string]bool{"window": status}
	json.NewEncoder(w).Encode(response)
}

func ToggleWindowHandler(w http.ResponseWriter, r *http.Request) {
	services.WindowMu.Lock()
	services.WindowOpen = !services.WindowOpen
	services.WindowEventFlag = !services.WindowEventFlag
	status := "closed"
	if services.WindowOpen {
		status = "open"
	}
	services.WindowMu.Unlock()
	fmt.Println("========== Toggled window. Now:", status, "==========")
	response := map[string]string{"window": status}
	json.NewEncoder(w).Encode(response)
}
