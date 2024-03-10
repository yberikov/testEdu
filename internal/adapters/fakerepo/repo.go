package fakerepo

import (
	"errors"
	"homework/internal/device"
	"sync"
)

type DeviceStorage struct {
	sync.Mutex
	devices map[string]device.Device
}

var (
	ErrNoSuchDevice        = errors.New("there is no such device")
	ErrDeviceAlreadyExists = errors.New("such device is already in database")
)

func NewDeviceStorage() *DeviceStorage {
	return &DeviceStorage{
		devices: make(map[string]device.Device),
	}
}

func (s *DeviceStorage) GetDeviceBySerialNum(serialNum string) (device.Device, error) {
	defer s.Unlock()
	s.Lock()
	if val, ok := s.devices[serialNum]; ok {
		return val, nil
	}
	return device.Device{}, ErrNoSuchDevice
}

func (s *DeviceStorage) CreateDevice(device device.Device) error {
	defer s.Unlock()
	s.Lock()
	if _, ok := s.devices[device.SerialNum]; ok {
		return ErrDeviceAlreadyExists
	}
	s.devices[device.SerialNum] = device
	return nil
}

func (s *DeviceStorage) DeleteDeviceBySerialNum(serialNum string) error {
	defer s.Unlock()
	s.Lock()
	if _, ok := s.devices[serialNum]; ok {
		delete(s.devices, serialNum)
		return nil
	}
	return ErrNoSuchDevice
}

func (s *DeviceStorage) UpdateDevice(device device.Device) error {
	defer s.Unlock()
	s.Lock()
	if _, ok := s.devices[device.SerialNum]; ok {
		s.devices[device.SerialNum] = device
		return nil
	}
	return ErrNoSuchDevice
}
