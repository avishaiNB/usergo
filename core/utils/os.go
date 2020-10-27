package utils

import "os"

// ProcessName ...
func ProcessName() string {
	name, _ := os.Executable()
	return name
}

// GetEnvVar ...
func GetEnvVar(key string) string {
	return os.Getenv(key)
}
