package usbdriver

import (
	"errors"

	"github.com/zserge/hid"
)

var (
	errDeviceNotFound = errors.New("device not found")
)

func GetDevice(vendorID, productID uint16) (hid.Device, error) {
	var dev hid.Device

	hid.UsbWalk(func(d hid.Device) {
		if d.Info().Product == 0x3008 && d.Info().Vendor == 0x1e71 {
			dev = d
		}
	})

	if dev == nil {
		return nil, errDeviceNotFound
	}

	return dev, nil
}
