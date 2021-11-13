package sensor

import (
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var (
	errDeviceNotFound = errors.New("sensors device not found")
	errParseError     = errors.New("sensors parse error")
)

func init() {
	// check if sensors is installed
	output, err := exec.Command("sensors", "-v").Output()

	if err != nil || !strings.Contains(string(output), "sensors version") {
		log.Fatal("sensors not installed")
	}

	log.Printf("sensors installed: %s", string(output))
}

func GetTemp(device string) (float64, error) {
	output, err := exec.Command("sensors").Output()

	if err != nil {
		return 0, err
	}

	outputString := string(output)

	if !strings.Contains(outputString, device) {
		return 0, errDeviceNotFound
	}

	// get correct device block

	afterDeviceBlock := strings.Split(outputString, device)

	if len(afterDeviceBlock) < 1 {
		return 0, errParseError
	}

	deviceBlock := strings.Split(afterDeviceBlock[1], "\n\n")[0]

	// get temp

	tempString := strings.Split(strings.Split(deviceBlock, "+")[1], "°C")[0]

	tempFloat, err := strconv.ParseFloat(tempString, 64)

	if err != nil {
		return 0, errParseError
	}

	return tempFloat, nil
}
