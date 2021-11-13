package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/yanando/cpu-temp/config"
	"github.com/yanando/cpu-temp/liquidctl"
	"github.com/yanando/cpu-temp/sensor"
)

type coefMatch struct {
	temp int
	coef float64
}

type Monitor struct {
	Config config.Config

	liquidctlDevice liquidctl.LiquidctlDevice

	// ordered from cool -> hot temps
	fanCurveCoefficients  []coefMatch
	pumpCurveCoefficients []coefMatch
}

func (m *Monitor) Start() {
	m.fanCurveCoefficients = calculateCoefficient(m.Config.FanCurve)
	m.pumpCurveCoefficients = calculateCoefficient(m.Config.PumpCurve)

	m.liquidctlDevice = liquidctl.LiquidctlDevice{
		ID: m.Config.LiquidctlDeviceID,
	}

	for {
		temp, err := sensor.GetTemp(m.Config.TemperatureDevice)

		if err != nil {
			log.Fatalf("Error getting temperature: %v", err)
		}

		fanspeed, pumpspeed := m.calculateSpeedsFromTemp(temp)

		fmt.Printf("Temperature: %.01f°C\nFan speed: %d%%\nPump speed: %d%%\n\n", temp, fanspeed, pumpspeed)

		err = m.liquidctlDevice.SetFanSpeed(fanspeed)

		if err != nil {
			log.Fatalf("Error setting fan speed: %v", err)
		}

		err = m.liquidctlDevice.SetPumpSpeed(pumpspeed)

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
