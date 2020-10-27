package configuration

import (
	viper "github.com/spf13/viper"
	tleutils "github.com/thelotter-enterprise/usergo/core/utils"
)

const (
	consulClientKey      = "CONSUL_CLIENT"
	defaultConsulAddress = "localhost:8500"
)

// Client ...
type Client interface {
	Get(key string, defaultValue interface{}) interface{}
}

type client struct {
	viperClient *viper.Viper
}

// NewConfiguration ...
func NewConfiguration() (Client, error) {
	viperClient := viper.New()
	viperClient.SetConfigType("json")
	err := addConsulProviderToClient(viperClient)

	if err != nil {
		return nil, err
	}

	return &client{
		viperClient: viperClient,
	}, nil
}

func addConsulProviderToClient(viperClient *viper.Viper) error {
	consulClient := getConsulAddress()
	viperClient.AddRemoteProvider("consul", consulClient, tleutils.ProcessName())
	err := viperClient.ReadRemoteConfig()

	if err != nil {
		return err
	}

	err = viperClient.WatchRemoteConfigOnChannel()

	if err != nil {
		return err
	}
	return nil
}

func getConsulAddress() string {
	consulClientAddress := tleutils.GetEnvVar(consulClientKey)
	if consulClientAddress == "" {
		return defaultConsulAddress
	}
	return consulClientAddress
}

func (configClient client) Get(key string, defaultValue interface{}) interface{} {
	value := configClient.getFromConsul(key)
	if value == nil {
		return defaultValue
	}

	return value
}

func (configClient client) getFromConsul(key string) interface{} {
	return configClient.viperClient.Get(key)
}
