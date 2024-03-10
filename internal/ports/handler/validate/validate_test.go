package validate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"homework/internal/device"
	"net"
	"regexp"
	"testing"
)

func TestValidateDevice(t *testing.T) {
	tests := []struct {
		name     string
		device   device.Device
		expected error
	}{
		{
			name:     "Valid Device",
			device:   device.Device{SerialNum: "12345", Model: "Model1", IP: "192.168.1.1"},
			expected: nil,
		},
		{
			name:     "Empty SerialNum",
			device:   device.Device{SerialNum: "", Model: "Model2", IP: "192.168.1.2"},
			expected: fmt.Errorf("field cannot be empty"),
		},
		{
			name:     "Empty Model",
			device:   device.Device{SerialNum: "67890", Model: "", IP: "192.168.1.3"},
			expected: fmt.Errorf("field cannot be empty"),
		},
		{
			name:     "Empty IP",
			device:   device.Device{SerialNum: "54321", Model: "Model3", IP: ""},
			expected: fmt.Errorf("field cannot be empty"),
		},
		{
			name:     "Invalid IP Format",
			device:   device.Device{SerialNum: "13579", Model: "Model4", IP: "invalidip"},
			expected: fmt.Errorf("IP field is in wrong format"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateDevice(test.device)
			assert.Equal(t, test.expected, err)
		})
	}
}

// using net.ParseIP to validate IP field of model, so fuzzing it with regex implementation
func FuzzIPCheck(f *testing.F) {
	testcases := []string{"192.0.2.1", "92.121.24.11", "192.61.4.77"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	f.Fuzz(func(t *testing.T, orig string) {
		if (net.ParseIP(orig) != nil) != re.MatchString(orig) {
			t.Logf("invalid ip adress: %q", orig)
		}
	})
}

func TestIsValidSerialNum(t *testing.T) {
	tests := []struct {
		serialNum string
		expected  error
	}{
		{"abc", nil},
		{"123", nil},
		{"a1b2c3", nil},
		{"ab", ErrSerialNumLength},
		{"123@", ErrSerialNumChar},
		{"", ErrSerialNumLength},
	}

	for _, test := range tests {
		t.Run(test.serialNum, func(t *testing.T) {
			err := IsValidSerialNum(test.serialNum)

			if (err != nil && test.expected == nil) || (err == nil && test.expected != nil) || (err != nil && test.expected != nil && err.Error() != test.expected.Error()) {
				t.Errorf("For serialNum '%s', got error '%v', expected '%v'", test.serialNum, err, test.expected)
			}
		})
	}
}
