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
		if d.Info().Product == productID && d.Info().Vendor == vendorID {
			dev = d
		}
	})

	if dev == nil {
		return nil, errDeviceNotFound
	}

	return dev, nil
}
