package config

import (
	"strings"
	"time"
)

func ParseTimeout(timeout string) (time.Duration, error) {
	if timeout == "" {
		return 0, nil
	}
	if !strings.HasSuffix(timeout, "s") && !strings.HasSuffix(timeout, "m") && !strings.HasSuffix(timeout, "h") {
		timeout += "s"
	}
	return time.ParseDuration(timeout)
}
