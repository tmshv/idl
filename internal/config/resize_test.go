package config

import (
	"testing"
)

func TestParseResize(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected [2]int
		err      bool
	}{
		{
			name:     "empty string",
			value:    "",
			expected: [2]int{0, 0},
			err:      false,
		},
		{
			name:     "valid resize",
			value:    "100x200",
			expected: [2]int{100, 200},
			err:      false,
		},
		{
			name:     "invalid format - missing x",
			value:    "100",
			expected: [2]int{0, 0},
			err:      true,
		},
		{
			name:     "invalid format - extra x",
			value:    "100x200x300",
			expected: [2]int{0, 0},
			err:      true,
		},
		{
			name:     "invalid width - not a number",
			value:    "abcx200",
			expected: [2]int{0, 0},
			err:      true,
		},
		{
			name:     "invalid height - not a number",
			value:    "100xdef",
			expected: [2]int{0, 0},
			err:      true,
		},
		{
			name:     "negative width",
			value:    "-100x200",
			expected: [2]int{-100, 200},
			err:      false, // Function doesn't explicitly validate positive numbers
		},
		{
			name:     "negative height",
			value:    "100x-200",
			expected: [2]int{100, -200},
			err:      false, // Function doesn't explicitly validate positive numbers
		},
		{
			name:     "large numbers",
			value:    "1234567890x9876543210",
			expected: [2]int{1234567890, 9876543210},
			err:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ParseResize(tc.value)
			if tc.err {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				if actual[0] != tc.expected[0] || actual[1] != tc.expected[1] {
					t.Errorf("expected %v, got %v", tc.expected, actual)
				}
			}
		})
	}
}
