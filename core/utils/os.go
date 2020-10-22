package utils

import "os"

// ProcessName ...
func ProcessName() string {
	name, _ := os.Executable()
	return name
}
