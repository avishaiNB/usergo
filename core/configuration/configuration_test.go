package configuration_test

import (
	"testing"

	"github.com/thelotter-enterprise/usergo/core/configuration"
)

func TestConfiguration(t *testing.T) {
	conf := configuration.NewNopConfigurationClient()
	stringValue := getStringValueFromConfig(conf, "Status")
	if stringValue != "Status" {
		t.Errorf("getStringValueFromConfig return wrong value %s ; want %s ;", stringValue, "Status")
	}

	interfaceValue := getInterfaceFromConfig(conf, "Status", 3)
	if interfaceValue.(int) != 3 {
		t.Errorf("getStringValueFromConfig return wrong value %v ; want %d ;", interfaceValue, 3)
	}
}

func getStringValueFromConfig(configClient configuration.Client, key string) string {
	return configClient.GetString(key)
}

func getInterfaceFromConfig(configClient configuration.Client, key string, defaultValue interface{}) interface{} {
	return configClient.Get(key, defaultValue)
}
