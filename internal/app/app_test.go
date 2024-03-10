package app

import (
	"github.com/stretchr/testify/mock"
	"homework/internal/adapters/fakerepo"
	"homework/internal/app/mocks"
	"homework/internal/device"
	"testing"
)

func TestCreateDevice(t *testing.T) {
	storageMock := mocks.NewDeviceStorage(t)
	service := NewService(storageMock)

	wantDevice := device.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	storageMock.On("CreateDevice", wantDevice).
		Return(nil)
	err := service.CreateDevice(wantDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	storageMock.On("GetDeviceBySerialNum", wantDevice.SerialNum).
		Return(wantDevice, nil)
	gotDevice, err := service.GetDevice(wantDevice.SerialNum)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if wantDevice != gotDevice {
		t.Errorf("want device %+#v not equal got %+#v", wantDevice, gotDevice)
	}
}

func TestCreateMultipleDevices(t *testing.T) {
	storageMock := mocks.NewDeviceStorage(t)
	service := NewService(storageMock)
	devices := []device.Device{
		{
			SerialNum: "123",
			Model:     "model1",
			IP:        "1.1.1.1",
		},
		{
			SerialNum: "124",
			Model:     "model2",
			IP:        "1.1.1.2",
		},
		{
			SerialNum: "125",
			Model:     "model3",
			IP:        "1.1.1.3",
		},
	}

	for _, d := range devices {
		storageMock.On("CreateDevice", d).
			Return(nil)
		err := service.CreateDevice(d)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	for _, wantDevice := range devices {
		storageMock.On("GetDeviceBySerialNum", wantDevice.SerialNum).
			Return(wantDevice, nil)
		gotDevice, err := service.GetDevice(wantDevice.SerialNum)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if wantDevice != gotDevice {
			t.Errorf("want device %+#v not equal got %+#v", wantDevice, gotDevice)
		}
	}
}

func TestCreateDuplicate(t *testing.T) {
	storageMock := mocks.NewDeviceStorage(t)
	service := NewService(storageMock)
	wantDevice := device.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	storageMock.On("CreateDevice", wantDevice).
		Return(nil).Once()
	err := service.CreateDevice(wantDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	storageMock.On("CreateDevice", wantDevice).
		Return(fakerepo.ErrDeviceAlreadyExists).Once()
	err = service.CreateDevice(wantDevice)
	if err == nil {
		t.Errorf("want error, but got nil")
	}

}

func TestGetDeviceUnexisting(t *testing.T) {
	storageMock := mocks.NewDeviceStorage(t)
	service := NewService(storageMock)
	wantDevice := device.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	storageMock.On("CreateDevice", wantDevice).
		Return(nil).Once()
	err := service.CreateDevice(wantDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	storageMock.On("GetDeviceBySerialNum", wantDevice.SerialNum).
		Return(wantDevice, nil).Maybe()
	storageMock.On("GetDeviceBySerialNum", mock.Anything).
		Return(device.Device{}, fakerepo.ErrNoSuchDevice)
	_, err = service.GetDevice("1")
	if err == nil {
		t.Error("want error, but got nil")
	}
}

func TestDeleteDevice(t *testing.T) {
	storageMock := mocks.NewDeviceStorage(t)
	service := NewService(storageMock)
	newDevice := device.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	storageMock.On("CreateDevice", newDevice).
		Return(nil).Once()

	err := service.CreateDevice(newDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	storageMock.On("DeleteDeviceBySerialNum", newDevice.SerialNum).
		Return(nil).Once()
	err = service.DeleteDevice(newDevice.SerialNum)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	storageMock.On("GetDeviceBySerialNum", newDevice.SerialNum).
		Return(device.Device{}, fakerepo.ErrNoSuchDevice)
	_, err = service.GetDevice(newDevice.SerialNum)
	if err == nil {
		t.Error("want error, but got nil")
	}
}

func TestDeleteDeviceUnexisting(t *testing.T) {
	storageMock := mocks.NewDeviceStorage(t)
	service := NewService(storageMock)

	storageMock.On("DeleteDeviceBySerialNum", mock.Anything).
		Return(fakerepo.ErrNoSuchDevice)
	err := service.DeleteDevice("123")
	if err == nil {
		t.Errorf("want error, but got nil")
	}
}

func TestUpdateDevice(t *testing.T) {
	storageMock := mocks.NewDeviceStorage(t)
	service := NewService(storageMock)
	testDevice := device.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	storageMock.On("CreateDevice", testDevice).
		Return(nil).Once()
	err := service.CreateDevice(testDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	newDevice := device.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.2",
	}
	storageMock.On("UpdateDevice", newDevice).
		Return(nil).Once()
	err = service.UpdateDevice(newDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	storageMock.On("GetDeviceBySerialNum", newDevice.SerialNum).
		Return(newDevice, nil)
	gotDevice, err := service.GetDevice(newDevice.SerialNum)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if gotDevice != newDevice {
		t.Errorf("new device %+#v not equal got device %+#v", newDevice, gotDevice)
	}
}

func TestUpdateDeviceUnexsting(t *testing.T) {
	storageMock := mocks.NewDeviceStorage(t)
	service := NewService(storageMock)
	testDevice := device.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	storageMock.On("CreateDevice", testDevice).
		Return(nil).Once()

	err := service.CreateDevice(testDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	newDevice := device.Device{
		SerialNum: "124",
		Model:     "model1",
		IP:        "1.1.1.2",
	}
	storageMock.On("UpdateDevice", newDevice).
		Return(fakerepo.ErrNoSuchDevice).Once()
	err = service.UpdateDevice(newDevice)
	if err == nil {
		t.Errorf("want err, but got nil")
	}
}
