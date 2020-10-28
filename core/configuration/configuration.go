package configuration

import (
	"time"

	viper "github.com/spf13/viper"
	tleutils "github.com/thelotter-enterprise/usergo/core/utils"
)

const (
	//JSON represents json config type
	JSON = "json"
	//YAML represents yaml config type
	YAML = "yaml"

	consulClientKey      = "CONSUL_CLIENT"
	defaultConsulAddress = "localhost:8500"
)

// Client represent contract of configuration client
type Client interface {
	//Get represents convention of get function with default value if cannot find key
	Get(string, interface{}) interface{}
	//GetString represents convention of get function which return string
	GetString(string) string
	//GetBool represents convention of get function which return bool
	GetBool(string) bool
	//GetInt represents convention of get function which return int
	GetInt(string) int
	//GetFloat represents convention of get function which return float
	GetFloat(string) float64
	//GetTime represents convention of get function which return time.Time
	GetTime(string) time.Time
	//GetDuration represents convention of get function which return time.Duration
	GetDuration(string) time.Duration
}

type client struct {
	viperClient *viper.Viper
}

// NewConfiguration create new configuration client , based on consul config f
// configType describe what is config file format (json , yaml etc.)
// viper client save all configs in inner cache
// viperclient monitor source file and pull new changes to cache
// Add 	_ "github.com/spf13/viper/remote" in main file for working with remote config source
func NewConfiguration(configType string) (Client, error) {
	viperClient := viper.New()
	viperClient.SetConfigType(configType)
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
	viperClient.AddRemoteProvider("consul", consulClient, "configurations/"+tleutils.ProcessName())
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

//Get return value from config , if cannot find key return default value
func (configClient client) Get(key string, defaultValue interface{}) interface{} {
	value := configClient.viperClient.Get(key)
	if value == nil {
		return defaultValue
	}

	return value
}

//GetString return value from config , if cannot find key return empty string
func (configClient client) GetString(key string) string {
	return configClient.viperClient.GetString(key)
}

//GetBool return value from config , if cannot find key return false
func (configClient client) GetBool(key string) bool {
	return configClient.viperClient.GetBool(key)
}

//GetInt return value from config , if cannot find key return 0
func (configClient client) GetInt(key string) int {
	return configClient.viperClient.GetInt(key)
}

//GetFloat return value from config , if cannot find key return 0.0
func (configClient client) GetFloat(key string) float64 {
	return configClient.viperClient.GetFloat64(key)
}

//GetTime return value from config , if cannot find key return 0001-01-01 00:00:00 +0000 UTC
func (configClient client) GetTime(key string) time.Time {
	return configClient.viperClient.GetTime(key)
}

//GetTime return value from config , if cannot find key return 0s
func (configClient client) GetDuration(key string) time.Duration {
	return configClient.viperClient.GetDuration(key)
}
