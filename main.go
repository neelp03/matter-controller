// main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/neelp03/matter-controller/handlers"
	"github.com/neelp03/matter-controller/services"
)

func isDeviceCommissioned() bool {
	cmd := exec.Command("../connectedhomeip/out/host/chip-tool", "operationalcredentials", "read", "fabrics", "1", "0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("!!!!!!!!!! Error checking commissioned status !!!!!!!!!!:", err)
		return false
	}

	outStr := string(output)
	if strings.Contains(outStr, "Fabrics: 0 entries") {
		return false
	}
	if strings.Contains(outStr, "Fabrics: 1 entries") && strings.Contains(outStr, "NodeID: 1") {
		return true
	}
	return false
}

func readTemperatureCelsius() (float64, error) {
	cmd := exec.Command("../connectedhomeip/out/host/chip-tool", "temperaturemeasurement", "read", "measured-value", "1", "1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("command failed: %v", err)
	}

	re := regexp.MustCompile(`(?m)MeasuredValue: (\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return 0, fmt.Errorf("!!!!!!!!!! could not parse temperature from output !!!!!!!!!!")
	}

	tempRaw, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("!!!!!!!!!! invalid temperature format !!!!!!!!!!: %v", err)
	}

	return float64(tempRaw) / 100.0, nil
}

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
		if err := services.PairDeviceOverBLE(ssid, password); err != nil {
			fmt.Println("!!!!!!!!!! Pairing failed !!!!!!!!!!:", err)
			return
		}
		fmt.Println("++++++++++ Pairing succeeded ++++++++++")
	}

	// Start periodic data backup
	backupInterval := 60 // seconds
	go services.Interval_backup(backupInterval)
	fmt.Println("========== Periodic data backup started every", backupInterval, "seconds ==========")

	// Start HTTP server
	http.HandleFunc("/temperature", handlers.TemperatureHandler)
	http.HandleFunc("/weather", handlers.WeatherHandler)

	fmt.Println("========== Server is now listening on port 8080 ==========")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println("!!!!!!!!!! Failed to start server !!!!!!!!!!:", err)
	}
}
