package utils

import "os"

// FileExists checks if a file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
