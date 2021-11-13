package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TemperatureDevice string `yaml:"temperatureDevice"`
	TemperaturePath   string `yaml:"temperaturePath"`

	ProductID uint16 `yaml:"productID"`
	VendorID  uint16 `yaml:"vendorID"`

	RefreshDelay float64 `yaml:"refreshDelay"`

	FanCurve  map[int]int `yaml:"fanCurve"`
	PumpCurve map[int]int `yaml:"pumpCurve"`
}

func init() {
	if _, err := os.Stat("config.yaml"); err == nil || !os.IsNotExist(err) {
		return
	}

	log.Println("config.yaml not found, generating")

	exampleConfig := Config{
		TemperatureDevice: "cpu_temp_device",
		TemperaturePath:   "/sys/devices/pci0000:00/0000:00:00.0/hwmon/hwmon2/temp1_input",

		// kraken Zx3 IDS
		ProductID: 0x3008,
		VendorID:  0x1e71,

		RefreshDelay: 1.1,
		FanCurve: map[int]int{
			35: 40,
			40: 45,
			50: 55,
			60: 75,
			70: 80,
			80: 100,
		},
		PumpCurve: map[int]int{
			35: 40,
			40: 45,
			50: 55,
			60: 75,
			70: 80,
			80: 100,
		},
	}

	configBytes, _ := yaml.Marshal(exampleConfig)

	err := os.WriteFile("config.yaml", configBytes, 0655)

	if err != nil {
		log.Fatalf("Error writing example config: %v", err)
	}
}

func GetConfig() Config {
	configBytes, err := os.ReadFile("config.yaml")

	if err != nil {
		log.Fatalf("Error opening config: %v", err)
	}

	var config Config

	err = yaml.Unmarshal(configBytes, &config)

	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	return config
}
