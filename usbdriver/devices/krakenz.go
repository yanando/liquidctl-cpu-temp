package devices

import (
	"fmt"
	"time"

	"github.com/zserge/hid"
)

const (
	readLength  = 64
	writeLength = 64
	ioTimeout   = time.Second * 10

	maxTemp = 60
)

type statusReport struct {
	LiquidTemperature float64

	FanSpeedRPM        int
	FanSpeedPercentage int

	PumpSpeedRPM        int
	PumpSpeedPercentage int
}

type KrakenZ struct {
	HidDev   hid.Device
	Firmware string
}

func (k *KrakenZ) Open() error {
	return k.HidDev.Open()
}

func (k *KrakenZ) Close() {
	k.HidDev.Close()
}

func (k *KrakenZ) Status() (statusReport, error) {
	_, err := k.Write(0x74, []byte{0x01})

	if err != nil {
		return statusReport{}, err
	}

	reportBytes, err := k.Read()

	if err != nil {
		return statusReport{}, err
	}

	// magic
	report := statusReport{
		LiquidTemperature:   float64(reportBytes[15]) + float64(reportBytes[16])/10,
		PumpSpeedRPM:        int(reportBytes[18])<<8 | int(reportBytes[17]),
		PumpSpeedPercentage: int(reportBytes[19]),
		FanSpeedRPM:         int(reportBytes[24])<<8 | int(reportBytes[23]),
		FanSpeedPercentage:  int(reportBytes[25]),
	}

	return report, nil
}

func (k *KrakenZ) Initialize() error {
	_, err := k.Write(0x10, []byte{0x01}) // firmware info

	if err != nil {
		return err
	}

	_, err = k.Write(0x20, []byte{0x03}) // lighting info

	if err != nil {
		return err
	}

	// update interval 1
	_, err = k.Write(0x70, []byte{0x02, 0x01, 0xb8, 0x01})

	if err != nil {
		return err
	}

	_, err = k.Write(0x70, []byte{0x01})

	if err != nil {
		return err
	}

	// -------------------------------------- led info?
	_, err = k.Read()

	if err != nil {
		return err
	}

	// -------------------------------------- firmware info

	res, err := k.Read()

	if err != nil {
		return err
	}

	k.Firmware = fmt.Sprintf("%d.%d.%d", res[17], res[18], res[19])

	fmt.Println(k.Firmware)

	// -------------------------------------- ??

	_, err = k.Read()

	if err != nil {
		return err
	}

	// -------------------------------------- ??

	_, err = k.Read()

	if err != nil {
		return err
	}

	return nil
}

func (k *KrakenZ) SetPumpSpeed(percentage int) error {
	payload := []byte{0x01, 0x00, 0x00}

	// set static speed for each temp from 21-59
	for temperature := 20; temperature < maxTemp; temperature++ {
		payload = append(payload, byte(percentage))
	}

	_, err := k.Write(0x72, payload)

	if err != nil {
		panic(err)
	}

	return nil
}

func (k *KrakenZ) SetFanSpeed(percentage int) error {
	payload := []byte{0x02, 0x00, 0x00}

	// set static speed for each temp from 21-59
	for temperature := 20; temperature < maxTemp; temperature++ {
		payload = append(payload, byte(percentage))
	}

	_, err := k.Write(0x72, payload)

	if err != nil {
		panic(err)
	}

	return nil
}

// Write writes to the kraken device with padding
func (k *KrakenZ) Write(reportID byte, data []byte) (int, error) {
	payload := []byte{reportID}

	padding := make([]byte, writeLength-len(data)-1)

	payload = append(payload, data...)
	payload = append(payload, padding...)

	return k.HidDev.Write(payload, ioTimeout)
}

func (k *KrakenZ) Read() ([]byte, error) {
	return k.HidDev.Read(readLength, ioTimeout)
}

func printPayload(payload []byte) {
	fmt.Print("[")

	for i, b := range payload {
		if i+1 == len(payload) {
			fmt.Printf("0x%x", b)
			break
		}
		fmt.Printf("0x%x, ", b)
	}

	fmt.Println("]")
}
