package services_test

import (
	"testing"

	"github.com/neelp03/matter-controller/services"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateRules(t *testing.T) {
	tests := []struct {
		name                string
		indoorTemp          float64
		outdoorTemp         float64
		rainRate            float64
		aqi                 float64
		expectedWindowState bool
	}{
		{
			name:                "Open window when indoor > maxComfort and (indoor - outdoor) >= deltaThreshold and safe conditions",
			indoorTemp:          78.0,
			outdoorTemp:         72.0,
			rainRate:            0.0,
			aqi:                 50.0,
			expectedWindowState: true,
		},
		{
			name:                "Close window when indoor < minComfort",
			indoorTemp:          65.0,
			outdoorTemp:         60.0,
			rainRate:            0.0,
			aqi:                 50.0,
			expectedWindowState: false,
		},
		{
			name:                "Close window when outdoor >= indoor",
			indoorTemp:          75.0,
			outdoorTemp:         76.0,
			rainRate:            0.0,
			aqi:                 50.0,
			expectedWindowState: false,
		},
		{
			name:                "Close window when rain > rainThreshold",
			indoorTemp:          75.0,
			outdoorTemp:         70.0,
			rainRate:            1.0,
			aqi:                 50.0,
			expectedWindowState: false,
		},
		{
			name:                "Close window when AQI > aqiThreshold",
			indoorTemp:          75.0,
			outdoorTemp:         70.0,
			rainRate:            0.0,
			aqi:                 150.0,
			expectedWindowState: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			services.ReadTemperature = func() (float64, error) {
				return tt.indoorTemp, nil
			}
			services.FetchOutdoorTemperature = func() (float64, error) {
				return tt.outdoorTemp, nil
			}
			services.FetchOutdoorRain = func() (float64, error) {
				return tt.rainRate, nil
			}
			services.FetchOutdoorAQI = func() (float64, error) {
				return tt.aqi, nil
			}
			services.OpenWindow = func() error {
				services.WindowMu.Lock()
				defer services.WindowMu.Unlock()
				services.WindowOpen = true
				return nil
			}
			services.CloseWindow = func() error {
				services.WindowMu.Lock()
				defer services.WindowMu.Unlock()
				services.WindowOpen = false
				return nil
			}

			services.WindowMu.Lock()
			services.WindowOpen = false
			services.WindowMu.Unlock()

			go services.EvaluateRules()
			time.Sleep(1 * time.Second)

			services.WindowMu.Lock()
			actualWindowState := services.WindowOpen
			services.WindowMu.Unlock()

			assert.Equal(t, tt.expectedWindowState, actualWindowState)
		})
	}
}