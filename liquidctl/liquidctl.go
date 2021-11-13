package liquidctl

import (
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var (
	errNotAPercentage = errors.New("not a percentage")
)

func init() {
	// check if liquidctl is installed
	output, err := exec.Command("liquidctl", "--version").Output()

	if err != nil || !strings.Contains(string(output), "liquidctl") {
		log.Fatal("liquidctl not installed")
	}

	log.Printf("liquidctl installed: %s", string(output))
}

func InitLiquidctl() {
	// initialize liquidctl
	err := exec.Command("liquidctl", "initialize", "all").Run()

	if err != nil {
		log.Fatalf("Error starting liquidctl: %v", err)
	}

	log.Println("liquidctl initialized")
}

type LiquidctlDevice struct {
	ID               string
	CurrentFanSpeed  int
	CurrentPumpSpeed int
}

func (l *LiquidctlDevice) SetFanSpeed(percentage int) error {
	// dont change speed if speed is the same
	if percentage == l.CurrentFanSpeed {
		return nil
	}

	if percentage < 0 || percentage > 100 {
		return errNotAPercentage
	}

	output, err := exec.Command("liquidctl", "--device="+l.ID, "set", "fan", "speed", strconv.Itoa(percentage)).Output()

	if err != nil {
		return err
	} else if strings.Contains(string(output), "Error") {
		return errors.New(string(output))
	}

	l.CurrentFanSpeed = percentage

	return nil
}

func (l *LiquidctlDevice) SetPumpSpeed(percentage int) error {
	// dont change speed if speed is the same
	if percentage == l.CurrentFanSpeed {
		return nil
	}

	if percentage < 0 || percentage > 100 {
		return errNotAPercentage
	}

	output, err := exec.Command("liquidctl", "--device="+l.ID, "set", "pump", "speed", strconv.Itoa(percentage)).Output()

	if err != nil {
		return err
	} else if strings.Contains(string(output), "Error") {
		return errors.New(string(output))
	}

	l.CurrentPumpSpeed = percentage

	return nil
}

func (l *LiquidctlDevice) GetSpeeds() (pumpSpeed int, fanSpeed int, err error) {
	output, err := exec.Command("liquidctl", "--device="+l.ID, "status").Output()

	if err != nil {
		return
	}

	outputString := string(output)

	pumpPercentage := strings.TrimSpace(strings.Split(strings.Split(outputString, "Pump duty")[1], "%")[0])
	fanPercentage := strings.TrimSpace(strings.Split(strings.Split(outputString, "Fan duty")[1], "%")[0])

	pumpSpeed, err = strconv.Atoi(pumpPercentage)

	if err != nil {
		return
	}

	fanSpeed, err = strconv.Atoi(fanPercentage)

	if err != nil {
		return
	}

	return
}
