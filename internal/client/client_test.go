package client_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	client2 "homework/internal/client"
	"homework/internal/config"
	"homework/internal/roundtripper"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoggingRoundTripper(t *testing.T) {
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(os.Stderr)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, Test Server!")
	}))
	defer testServer.Close()

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: roundtripper.LoggingRoundTripper{Next: http.DefaultTransport},
	}

	resp, err := client.Get(testServer.URL)
	if err != nil {
		t.Fatalf("Error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	expectedOutgoingLog := fmt.Sprintf("Outgoing Request: GET %s map[]", testServer.URL)
	expectedIncomingLog := fmt.Sprintf("Incoming Response")

	assert.Contains(t, logBuffer.String(), expectedOutgoingLog)
	assert.Contains(t, logBuffer.String(), expectedIncomingLog)
}

func TestBreakerRoundTripper(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/success" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Success!")
		} else if r.URL.Path == "/server-error" {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal Server Error")
		}
	}))
	defer testServer.Close()

	cb := roundtripper.CreateCb()
	breakerRoundTripper := roundtripper.BreakerRoundTripper{Cb: cb, Next: http.DefaultTransport}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: breakerRoundTripper,
	}

	/// Test a successful request
	resp, err := client.Get(testServer.URL + "/success")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test a failing request to trigger the circuit breaker
	resp, err = client.Get(testServer.URL + "/server-error")
	assert.NotNil(t, err)

	// Extract the error message without the appended URL
	expectedErrorMessage := "Server error: 500 Internal Server Error"
	actualErrorMessage := err.Error()
	assert.Contains(t, actualErrorMessage, expectedErrorMessage)
	assert.Nil(t, resp)

	time.Sleep(6 * time.Second)

	resp, err = client.Get(testServer.URL + "/success")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserRoundTripper(t *testing.T) {
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current directory: %v", err)
	}

	configPath := filepath.Join(currentDir, "..", "..", "config", "local.yaml")

	cfg := config.LoadConfig(configPath)

	fakeUserRoundTripper := &FakeRoundTripper{}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	client := client2.NewClient(cfg, fakeUserRoundTripper)

	resp, err := client.Get(testServer.URL)
	assert.Nil(t, err)

	// Ensure that the fake UserRoundTripper was called
	assert.True(t, fakeUserRoundTripper.Called)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// FakeRoundTripper is a fake implementation of http.RoundTripper for testing
type FakeRoundTripper struct {
	Called  bool
	Request *http.Request
}

func (f *FakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	f.Called = true
	f.Request = req

	// Return a dummy response
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       nil,
	}, nil
}
