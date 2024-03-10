package handler

import (
	"bytes"
	"encoding/json"
	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework/internal/adapters/fakerepo"
	"homework/internal/app"
	"homework/internal/device"
	"homework/internal/ports/handler/mocks"
	"homework/internal/ports/handler/validate"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_handleGetDevice(t *testing.T) {

	tests := []struct {
		name          string
		method        string
		serialNum     string
		expectedCode  int
		expectedError error
		respErr       error
	}{
		{
			name:         "Success",
			method:       "GET",
			serialNum:    "1234",
			expectedCode: http.StatusOK,
		},
		{
			name:          "No Such Device",
			method:        "GET",
			serialNum:     "1234",
			expectedCode:  http.StatusBadRequest,
			expectedError: fakerepo.ErrNoSuchDevice,
			respErr:       fakerepo.ErrNoSuchDevice,
		},
		{
			name:          "Invalid http Method",
			serialNum:     "1234",
			method:        "POST",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidMethod,
		},
		{
			name:          "Invalid SerialNum",
			method:        "GET",
			serialNum:     "invalid@serial",
			expectedCode:  http.StatusBadRequest,
			expectedError: validate.ErrSerialNumChar,
		},
		{
			name:          "Invalid SerialNum",
			method:        "GET",
			serialNum:     "i",
			expectedCode:  http.StatusBadRequest,
			expectedError: validate.ErrSerialNumLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceMock := mocks.NewService(t)
			h := &Handler{
				service: serviceMock,
			}
			handler := h.InitRoutes()
			if tt.respErr == nil {
				serviceMock.On("GetDevice", tt.serialNum).
					Return(device.Device{SerialNum: tt.serialNum, Model: "HP", IP: "111.111.111.111"}, nil).Maybe()
			} else {
				serviceMock.On("GetDevice", tt.serialNum).
					Return(device.Device{}, tt.respErr).Maybe()
			}
			req, err := http.NewRequest(tt.method, "/getDevice", bytes.NewReader([]byte{}))
			require.NoError(t, err)

			req.Header.Set("serialNum", tt.serialNum)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedCode == http.StatusBadRequest {
				expectedError := MyError{Message: tt.expectedError.Error()}
				actualError := MyError{}

				err := json.Unmarshal(rr.Body.Bytes(), &actualError)
				if err != nil {
					assert.Fail(t, "error of unmarshalling error")
				}

				assert.Equal(t, expectedError, actualError)
			}
		})
	}
}

func TestHandler_handleDeleteDevice(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		serialNum     string
		expectedCode  int
		expectedError error
		respErr       error
	}{
		{
			name:         "Success",
			method:       "DELETE",
			serialNum:    "1234",
			expectedCode: http.StatusOK,
		},
		{
			name:          "No Such Device",
			method:        "DELETE",
			serialNum:     "1234",
			expectedCode:  http.StatusBadRequest,
			expectedError: fakerepo.ErrNoSuchDevice,
			respErr:       fakerepo.ErrNoSuchDevice,
		},
		{
			name:          "Invalid http Method",
			serialNum:     "1234",
			method:        "POST",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidMethod,
		},
		{
			name:          "Invalid SerialNum: char",
			method:        "DELETE",
			serialNum:     "invalid@serial",
			expectedCode:  http.StatusBadRequest,
			expectedError: validate.ErrSerialNumChar,
		},
		{
			name:          "Invalid SerialNum: len",
			method:        "DELETE",
			serialNum:     "i",
			expectedCode:  http.StatusBadRequest,
			expectedError: validate.ErrSerialNumLength,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceMock := mocks.NewService(t)
			h := &Handler{
				service: serviceMock,
			}
			handler := h.InitRoutes()
			if tt.respErr == nil {
				serviceMock.On("DeleteDevice", tt.serialNum).
					Return(nil).Maybe()
			} else {
				serviceMock.On("DeleteDevice", tt.serialNum).
					Return(tt.respErr).Maybe()
			}
			req, err := http.NewRequest(tt.method, "/deleteDevice", bytes.NewReader([]byte{}))
			require.NoError(t, err)

			req.Header.Set("serialNum", tt.serialNum)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedCode == http.StatusBadRequest {
				expectedError := MyError{Message: tt.expectedError.Error()}
				actualError := MyError{}

				err := json.Unmarshal(rr.Body.Bytes(), &actualError)
				if err != nil {
					assert.Fail(t, "error of unmarshalling error")
				}

				assert.Equal(t, expectedError, actualError)
			}
		})
	}
}

func TestHandler_handleCreateDevice(t *testing.T) {

	tests := []struct {
		name          string
		method        string
		serialNum     string
		model         string
		ip            string
		expectedCode  int
		expectedError error
		respErr       error
	}{
		{
			name:         "Success",
			method:       "POST",
			serialNum:    "1234",
			model:        "hp",
			ip:           "121.121.212.121",
			expectedCode: http.StatusOK,
		},
		{
			name:          "respErr",
			method:        "POST",
			serialNum:     "1234",
			model:         "hp",
			ip:            "121.121.212.121",
			expectedCode:  http.StatusBadRequest,
			expectedError: fakerepo.ErrNoSuchDevice,
			respErr:       fakerepo.ErrNoSuchDevice,
		},
		{
			name:          "Invalid http Method",
			method:        "GET",
			serialNum:     "1234",
			model:         "hp",
			ip:            "121.121.212.121",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidMethod,
		},
		{
			name:          "Invalid device: empty field",
			method:        "POST",
			serialNum:     "",
			model:         "hp",
			ip:            "121.121.212.121",
			expectedCode:  http.StatusBadRequest,
			expectedError: validate.ErrDeviceEmptyField,
		},
		{
			name:          "Invalid device: IP",
			method:        "POST",
			serialNum:     "3213",
			model:         "hp",
			ip:            "121.121ad.212.121.1311",
			expectedCode:  http.StatusBadRequest,
			expectedError: validate.ErrDeviceInvalidIP,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceMock := mocks.NewService(t)
			h := &Handler{
				service: serviceMock,
			}
			handler := h.InitRoutes()
			if tt.respErr == nil {
				serviceMock.On("CreateDevice", device.Device{SerialNum: tt.serialNum, Model: tt.model, IP: tt.ip}).
					Return(nil).Maybe()
			} else {
				serviceMock.On("CreateDevice", device.Device{SerialNum: tt.serialNum, Model: tt.model, IP: tt.ip}).
					Return(tt.respErr).Maybe()
			}
			req, err := http.NewRequest(tt.method, "/createDevice", bytes.NewReader([]byte{}))
			require.NoError(t, err)

			req.Header.Set("serialNum", tt.serialNum)
			req.Header.Set("Model", tt.model)
			req.Header.Set("IP", tt.ip)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedCode == http.StatusBadRequest {
				expectedError := MyError{Message: tt.expectedError.Error()}
				actualError := MyError{}

				err := json.Unmarshal(rr.Body.Bytes(), &actualError)
				if err != nil {
					assert.Fail(t, "error of unmarshalling error")
				}

				assert.Equal(t, expectedError, actualError)
			}
		})
	}
}

func TestHandler_handleUpdateDevice(t *testing.T) {

	tests := []struct {
		name          string
		method        string
		serialNum     string
		model         string
		ip            string
		expectedCode  int
		expectedError error
		respErr       error
	}{
		{
			name:         "Success",
			method:       "PUT",
			serialNum:    "1234",
			model:        "hp",
			ip:           "121.121.212.121",
			expectedCode: http.StatusOK,
		},
		{
			name:          "respErr",
			method:        "PUT",
			serialNum:     "1234",
			model:         "hp",
			ip:            "121.121.212.121",
			expectedCode:  http.StatusBadRequest,
			respErr:       fakerepo.ErrNoSuchDevice,
			expectedError: fakerepo.ErrNoSuchDevice,
		},
		{
			name:          "Invalid http Method",
			method:        "GET",
			serialNum:     "1234",
			model:         "hp",
			ip:            "121.121.212.121",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidMethod,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceMock := mocks.NewService(t)
			h := &Handler{
				service: serviceMock,
			}
			handler := h.InitRoutes()
			if tt.respErr == nil {
				serviceMock.On("UpdateDevice", device.Device{SerialNum: tt.serialNum, Model: tt.model, IP: tt.ip}).
					Return(nil).Maybe()
			} else {
				serviceMock.On("UpdateDevice", device.Device{SerialNum: tt.serialNum, Model: tt.model, IP: tt.ip}).
					Return(tt.respErr).Maybe()
			}
			req, err := http.NewRequest(tt.method, "/updateDevice", bytes.NewReader([]byte{}))
			req.Header.Set("serialNum", tt.serialNum)
			req.Header.Set("Model", tt.model)
			req.Header.Set("IP", tt.ip)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code)
			if tt.expectedCode == http.StatusBadRequest {
				expectedError := MyError{Message: tt.expectedError.Error()}
				actualError := MyError{}

				err := json.Unmarshal(rr.Body.Bytes(), &actualError)
				if err != nil {
					assert.Fail(t, "error of unmarshalling error")
				}

				assert.Equal(t, expectedError, actualError)
			}
		})
	}
}

type FakerModel struct {
	SerialNum string `faker:"word"`
	Model     string `faker:"word"`
	IP        string `faker:"ipv4"`
}

func BenchmarkCreateDeviceHandle(b *testing.B) {
	b.Run("Endpoint: /createDevice", func(b *testing.B) {

		h := &Handler{
			service: app.NewService(fakerepo.NewDeviceStorage()),
		}
		handler := h.InitRoutes()

		req, _ := http.NewRequest(
			"PUT", "/createDevice", nil)

		w := httptest.NewRecorder()

		// Turn on memory stats
		b.ReportAllocs()
		b.ResetTimer()

		// Execute the handler
		// with a request, `b.N` times
		for i := 0; i < b.N; i++ {
			a := FakerModel{}
			err := faker.FakeData(&a)
			if err != nil {
				continue
			}
			req.Header.Set("SerialNum", a.SerialNum)
			req.Header.Set("Model", a.Model)
			req.Header.Set("IP", a.IP)

			handler.ServeHTTP(w, req)
		}
	})
}
