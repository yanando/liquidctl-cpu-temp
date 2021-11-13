package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/yanando/cpu-temp/config"
	"github.com/yanando/cpu-temp/sensor"
	"github.com/yanando/cpu-temp/usbdriver"
	"github.com/yanando/cpu-temp/usbdriver/devices"
)

type coefMatch struct {
	temp int
	coef float64
}

type Monitor struct {
	Config config.Config

	// ordered from cool -> hot temps
	fanCurveCoefficients  []coefMatch
	pumpCurveCoefficients []coefMatch
}

func (m *Monitor) Start() {
	m.fanCurveCoefficients = calculateCoefficient(m.Config.FanCurve)
	m.pumpCurveCoefficients = calculateCoefficient(m.Config.PumpCurve)

	device, err := usbdriver.GetDevice(m.Config.VendorID, m.Config.ProductID)

	if err != nil {
		log.Fatalf("Error getting device: %v", err)
	}

	kraken := devices.KrakenZ{HidDev: device}

	err = kraken.Open()

	if err != nil {
		log.Fatalf("Error opening krakenz device")
	}

	defer kraken.Close()
	for {
		var temp float64

		if m.Config.TemperaturePath != "" {
			t, err := sensor.GetKernelTemp(m.Config.TemperaturePath)

			if err != nil {
				log.Fatalf("Error getting temperature: %v", err)
			}

			temp = t
		} else {
			t, err := sensor.GetTemp(m.Config.TemperatureDevice)

			if err != nil {
				log.Fatalf("Error getting temperature: %v", err)
			}

			temp = t
		}

		fanspeed, pumpspeed := m.calculateSpeedsFromTemp(temp)

		fmt.Printf("Temperature: %.01f°C\nFan speed: %d%%\nPump speed: %d%%\n\n", temp, fanspeed, pumpspeed)

		err := kraken.SetFanSpeed(fanspeed)

		if err != nil {
			log.Fatalf("Error setting fan speed: %v", err)
		}

		err = kraken.SetPumpSpeed(pumpspeed)

		if err != nil {
			log.Fatalf("Error setting pump speed: %v", err)
		}

		time.Sleep(time.Second * time.Duration(m.Config.RefreshDelay))
	}
}

// calculateSpeedsFromTemp returns (fanspeed, pumpspeed)
func (m *Monitor) calculateSpeedsFromTemp(temp float64) (int, int) {

	fanSpeed := calculateSpeedsFromTemp(temp, m.fanCurveCoefficients, m.Config.FanCurve)
	pumpSpeed := calculateSpeedsFromTemp(temp, m.pumpCurveCoefficients, m.Config.PumpCurve)

	return fanSpeed, pumpSpeed
}
