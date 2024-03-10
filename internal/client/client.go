package client

import (
	"homework/internal/config"
	"homework/internal/roundtripper"
	"net/http"
)

func NewClient(cfg *config.Config, customRoundTrippers ...http.RoundTripper) *http.Client {
	transport := http.DefaultTransport

	for _, rt := range customRoundTrippers {
		transport = roundtripper.UserRoundTripper{Transport: transport, UserRoundTrip: rt}
	}

	transport = roundtripper.BreakerRoundTripper{Cb: roundtripper.CreateCb(), Next: transport}

	transport = roundtripper.LoggingRoundTripper{Next: transport}
	return &http.Client{
		Transport: transport,
		Timeout:   cfg.HttpClient.Timeout,
	}
}
