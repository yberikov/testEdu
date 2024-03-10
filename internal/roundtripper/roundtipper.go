package roundtripper

import (
	"fmt"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"time"
)

type LoggingRoundTripper struct {
	Next http.RoundTripper
}

func (l LoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Printf("Outgoing Request: %s %s %v", r.Method, r.URL, r.Header)

	resp, err := l.Next.RoundTrip(r)

	if resp != nil {
		log.Printf("Incoming Response: %d %v", resp.StatusCode, resp.Header)
	} else if err != nil {
		log.Printf("Error in Response: %v", err)
	}

	return resp, err
}

type BreakerRoundTripper struct {
	Cb   *gobreaker.CircuitBreaker
	Next http.RoundTripper
}

func (l BreakerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {

	resp, err := l.Cb.Execute(func() (any, error) {
		resp, err := l.Next.RoundTrip(r)
		if err != nil {
			log.Printf("Error making request: %v", err)
		} else if resp.StatusCode == 500 {
			err = fmt.Errorf("Server error: %s", resp.Status)
		}
		return resp, err
	})

	if err != nil {
		return nil, err
	}

	if httpResponse, ok := resp.(*http.Response); ok {
		return httpResponse, nil
	}
	return nil, fmt.Errorf("Unexpected response type")
}

func CreateCb() *gobreaker.CircuitBreaker {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "HTTP GET",
		Timeout: 5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	})
	return cb
}

type UserRoundTripper struct {
	Transport     http.RoundTripper
	UserRoundTrip http.RoundTripper
}

func (urt UserRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return urt.UserRoundTrip.RoundTrip(req)
}
