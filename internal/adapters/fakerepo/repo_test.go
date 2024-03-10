package fakerepo

import (
	"github.com/stretchr/testify/suite"
	"homework/internal/device"
	"log"
	"sync"
	"testing"
)

type MyTestSuite struct {
	suite.Suite
	devices map[string]device.Device
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}

// Create some initial values in our map before tests
func (s *MyTestSuite) SetupTest() {
	log.Println("SetupSuite()")
	s.devices = map[string]device.Device{"1235": {SerialNum: "1235", Model: "HP", IP: "121.121.121.121"},
		"7112": {SerialNum: "1235", Model: "HP", IP: "121.121.121.121"}}
}

func (s *MyTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")

	// Imitation of deletion of created database ( if it used)
	s.devices = map[string]device.Device{}
}

func (s *MyTestSuite) TestDeviceStorage_GetDeviceBySerialNumDevice() {
	tests := []struct {
		name      string
		serialNum string
		err       error
	}{
		{
			name:      "success",
			serialNum: "1235",
			err:       nil,
		},
		{
			name:      "noSuchDevice",
			serialNum: "6143",
			err:       ErrNoSuchDevice,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			storage := &DeviceStorage{
				Mutex:   sync.Mutex{},
				devices: s.devices,
			}
			got, err := storage.GetDeviceBySerialNum(tt.serialNum)
			if err != tt.err {
				s.T().Error("Expected and result error is not equal")
			}
			if got != s.devices[tt.serialNum] {
				s.T().Error("Expected and result device is not equal")
			}
		})
	}
}

func (s *MyTestSuite) TestDeviceStorage_CreateDevice() {
	tests := []struct {
		name   string
		device device.Device
		err    error
	}{
		{
			name:   "success",
			device: device.Device{SerialNum: "4444", Model: "KOP", IP: "013.33.121.6"},
			err:    nil,
		},
		{
			name:   "DeviceAlreadyExists",
			device: device.Device{SerialNum: "1235", Model: "HP", IP: "121.121.121.121"},
			err:    ErrDeviceAlreadyExists,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			storage := &DeviceStorage{
				Mutex:   sync.Mutex{},
				devices: s.devices,
			}
			err := storage.CreateDevice(tt.device)
			if err != tt.err {
				s.T().Error("Expected and result error is not equal")
			}
		})
	}
}

func (s *MyTestSuite) TestDeviceStorage_DeleteDeviceBySerialNumDevice() {
	tests := []struct {
		name      string
		serialNum string
		err       error
	}{
		{
			name:      "success",
			serialNum: "1235",
			err:       nil,
		},
		{
			name:      "noSuchDevice",
			serialNum: "6143",
			err:       ErrNoSuchDevice,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			storage := &DeviceStorage{
				Mutex:   sync.Mutex{},
				devices: s.devices,
			}
			err := storage.DeleteDeviceBySerialNum(tt.serialNum)
			if err != tt.err {
				s.T().Error("Expected and result error is not equal")
			}
		})
	}
}

func (s *MyTestSuite) TestDeviceStorage_UpdateDevice() {
	tests := []struct {
		name   string
		device device.Device
		err    error
	}{
		{
			name:   "success",
			device: device.Device{SerialNum: "1235", Model: "ASUS", IP: "1.1.4.44"},
			err:    nil,
		},
		{
			name:   "NoSuchDevice",
			device: device.Device{SerialNum: "4444", Model: "KOP", IP: "013.33.121.6"},
			err:    ErrNoSuchDevice,
		},
	}
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			storage := &DeviceStorage{
				Mutex:   sync.Mutex{},
				devices: s.devices,
			}
			err := storage.UpdateDevice(tt.device)
			if err != tt.err {
				s.T().Error("Expected and result error is not equal")
			}
		})
	}
}
