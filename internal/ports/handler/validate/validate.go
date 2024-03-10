package validate

import (
	"errors"
	"homework/internal/device"
	"net"
	"regexp"
)

var (
	rexegpSerialNum     = "^[0-9a-zA-Z]+$"
	ErrSerialNumLength  = errors.New("serialNum should be at least 3 characters long")
	ErrSerialNumChar    = errors.New("serialNum should contain only digits or letters")
	ErrDeviceEmptyField = errors.New("field cannot be empty")
	ErrDeviceInvalidIP  = errors.New("IP field is in wrong format")
)

func ValidateDevice(device device.Device) error {
	if device.SerialNum == "" || device.Model == "" || device.IP == "" {
		return ErrDeviceEmptyField
	}
	if err := IsValidSerialNum(device.SerialNum); err != nil {
		return err
	}
	if net.ParseIP(device.IP) == nil {
		return ErrDeviceInvalidIP
	}
	return nil
}

// IsValidSerialNum should be at least 3 characters in length and contain only digits or letters.
func IsValidSerialNum(serialNum string) error {
	if len(serialNum) < 3 {
		return ErrSerialNumLength
	}

	if !(regexp.MustCompile(rexegpSerialNum).MatchString(serialNum)) {
		return ErrSerialNumChar
	}
	return nil

}
