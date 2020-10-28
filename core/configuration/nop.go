package configuration

import "time"

type nopConfigurationClient struct{}

// NewNopConfigurationClient returns a log.Logger that doesn't do anything.
// Should be used `for testing only
func NewNopConfigurationClient() Client {
	return nopConfigurationClient{}
}

func (nopConfigurationClient nopConfigurationClient) Get(key string, defaultValue interface{}) interface{} {
	if key == "1" {
		return "1"
	}
	return defaultValue
}

func (nopConfigurationClient nopConfigurationClient) GetString(key string) string {
	return key
}

func (nopConfigurationClient nopConfigurationClient) GetBool(key string) bool {
	if key == "1" {
		return true
	}
	return false
}

func (nopConfigurationClient nopConfigurationClient) GetInt(key string) int {
	if key == "1" {
		return 2
	}
	return 1
}

func (nopConfigurationClient nopConfigurationClient) GetFloat(key string) float64 {
	if key == "1" {
		return 1.3
	}
	return 1
}

func (nopConfigurationClient nopConfigurationClient) GetTime(key string) time.Time {
	return time.Now()
}

func (nopConfigurationClient nopConfigurationClient) GetDuration(key string) time.Duration {
	return time.Duration(1)
}
