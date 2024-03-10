package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"homework/internal/device"
	"homework/internal/ports/handler/validate"
	"log"
	"net/http"
)

var (
	ErrInvalidMethod = errors.New("invalid http method")
)

type MyError struct {
	Message string `json:"message"`
}

func (e *MyError) Error() string {
	return e.Message
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	myErr := &MyError{Message: err.Error()}
	jsonErr, err := json.Marshal(myErr)
	if err != nil {
		http.Error(w, "Failed to marshal error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonErr)
	if err != nil {
		log.Printf("Failed to write error, %d", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) handleGetDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		writeError(w, http.StatusBadRequest, ErrInvalidMethod)
		return
	}
	serialNum := r.Header.Get("serialNum")
	if err := validate.IsValidSerialNum(serialNum); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	d, err := h.service.GetDevice(serialNum)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	fmt.Printf("Device with SerialNum %s, Model %s, IP %s\n", d.SerialNum, d.Model, d.IP)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleCreateDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeError(w, http.StatusBadRequest, ErrInvalidMethod)
		return
	}
	serialNum := r.Header.Get("serialNum")
	Model := r.Header.Get("Model")
	IP := r.Header.Get("IP")
	device := device.Device{SerialNum: serialNum, Model: Model, IP: IP}
	if err := validate.ValidateDevice(device); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	err := h.service.CreateDevice(device)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("Device succsesfully created")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleDeleteDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		writeError(w, http.StatusBadRequest, ErrInvalidMethod)
		return
	}
	serialNum := r.Header.Get("serialNum")
	if err := validate.IsValidSerialNum(serialNum); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	err := h.service.DeleteDevice(serialNum)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("Device succsesfully deleted")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleUpdateDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		writeError(w, http.StatusBadRequest, ErrInvalidMethod)
		return
	}
	serialNum := r.Header.Get("serialNum")
	Model := r.Header.Get("Model")
	IP := r.Header.Get("IP")
	device := device.Device{SerialNum: serialNum, Model: Model, IP: IP}
	if err := validate.ValidateDevice(device); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	err := h.service.UpdateDevice(device)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("Device successfully updated")
	w.WriteHeader(http.StatusOK)
}
