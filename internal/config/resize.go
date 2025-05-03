package config

import (
	"errors"
	"strconv"
	"strings"
)

// ParseResize parses a string representing a resize configuration.
// Example: "100x200" -> [100, 200]
func ParseResize(value string, min [2]int, max [2]int) ([2]int, error) {
	if value == "" {
		return [2]int{0, 0}, nil
	}

	parts := strings.Split(value, "x")
	if len(parts) != 2 {
		return [2]int{0, 0}, errors.New("invalid resize format")
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return [2]int{0, 0}, errors.New("invalid width")
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return [2]int{0, 0}, errors.New("invalid height")
	}

	if width < min[0] || height < min[1] {
		return [2]int{0, 0}, errors.New("size too small")
	}
	if width > max[0] || height > max[1] {
		return [2]int{0, 0}, errors.New("size too large")
	}

	return [2]int{width, height}, nil
}
