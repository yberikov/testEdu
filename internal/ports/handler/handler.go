package handler

import (
	"homework/internal/device"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.36.0 --name=Service
type Service interface {
	GetDevice(string) (device.Device, error)
	CreateDevice(device.Device) error
	DeleteDevice(string) error
	UpdateDevice(device.Device) error
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/getDevice", h.handleGetDevice)
	mux.HandleFunc("/createDevice", h.handleCreateDevice)
	mux.HandleFunc("/deleteDevice", h.handleDeleteDevice)
	mux.HandleFunc("/updateDevice", h.handleUpdateDevice)
	return mux
}
