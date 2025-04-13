package services

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"github.com/neelp03/matter-controller/utils"
)

func IsDeviceCommissioned() bool {
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

func PairDeviceOverBLE(ssid, password string) error {
	fmt.Println("========== Device not commissioned. Attempting to pair via BLE... ==========")
	cmd := exec.Command("../connectedhomeip/out/host/chip-tool", "pairing", "ble-wifi", "1", ssid, password, "20202021", "3840")
	cmdOut, err := cmd.CombinedOutput()
	fmt.Println("========== Pairing output ==========")
	fmt.Println(string(cmdOut))
	if err != nil {
		return fmt.Errorf("!!!!!!!!!! Pairing failed !!!!!!!!!!: %v", err)
	}
	fmt.Println("++++++++++ Pairing succeeded ++++++++++")
	return nil
}

func ReadTemperature() (float64, error) {
	cmd := exec.Command("../connectedhomeip/out/host/chip-tool", "temperaturemeasurement", "read", "measured-value", "1", "1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("!!!!!!!!!! command failed !!!!!!!!!!: %v", err)
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

	celsius := float64(tempRaw) / 100.0
	fahrenheit := utils.CToF(celsius)

	return fahrenheit, nil
}

