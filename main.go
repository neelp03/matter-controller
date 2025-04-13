// main.go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/neelp03/matter-controller/handlers"
	"github.com/neelp03/matter-controller/services"
	"github.com/neelp03/matter-controller/utils"
)

func main() {
	fmt.Println("========== Starting Matter controller server ==========")

	utils.LoadEnv()
	ssid := os.Getenv("WIFI_SSID")
	password := os.Getenv("WIFI_PASSWORD")
	if ssid == "" || password == "" {
		fmt.Println("!!!!!!!!!! WIFI_SSID and WIFI_PASSWORD must be set in the environment !!!!!!!!!!")
		return
	}

	if services.IsDeviceCommissioned() {
		fmt.Println("++++++++++ Device is already commissioned ++++++++++")
	} else {
		fmt.Println("========== Device not commissioned. Attempting to pair via BLE... ==========")
		if err := services.PairDevice(ssid, password); err != nil {
			fmt.Println("!!!!!!!!!! Pairing failed !!!!!!!!!!:", err)
			return
		}
		fmt.Println("++++++++++ Pairing succeeded ++++++++++")
	}

	http.HandleFunc("/temperature", handlers.TemperatureHandler)
	http.HandleFunc("/weather", handlers.WeatherHandler)

	fmt.Println("========== Server is now listening on port 8080 ==========")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println("!!!!!!!!!! Failed to start server !!!!!!!!!!:", err)
	}
}