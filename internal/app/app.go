package app

import (
	"homework/internal/device"
)

//go:generate go run github.com/vektra/mockery/v2@v2.36.0 --name=DeviceStorage
type DeviceStorage interface {
	GetDeviceBySerialNum(serialNum string) (device.Device, error)
	CreateDevice(device device.Device) error
	DeleteDeviceBySerialNum(serialNum string) error
	UpdateDevice(device device.Device) error
}

type DeviceService struct {
	storage DeviceStorage
}

func NewService(storage DeviceStorage) *DeviceService {
	return &DeviceService{
		storage: storage,
	}
}

func (s *DeviceService) GetDevice(serialNum string) (device.Device, error) {
	d, err := s.storage.GetDeviceBySerialNum(serialNum)
	if err != nil {
		return device.Device{}, err
	}
	return d, nil
}

func (s *DeviceService) CreateDevice(device device.Device) error {
	if err := s.storage.CreateDevice(device); err != nil {
		return err
	}
	return nil
}

func (s *DeviceService) DeleteDevice(serialNum string) error {
	if err := s.storage.DeleteDeviceBySerialNum(serialNum); err != nil {
		return err
	}
	return nil
}

func (s *DeviceService) UpdateDevice(device device.Device) error {

	err := s.storage.UpdateDevice(device)
	if err != nil {
		return err
	}
	return nil
}
