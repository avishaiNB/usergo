package client

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
)

// ServiceDiscovery ...
type ServiceDiscovery struct {
	ConsulAPI    *consulapi.Client
	ConsulClient consul.Client
	Logger       log.Logger
}

// NewServiceDiscovery ...
func NewServiceDiscovery(logger log.Logger, consulAddress string) (ServiceDiscovery, error) {
	var err error
	sd := ServiceDiscovery{
		Logger: logger,
	}
	config := &consulapi.Config{
		Address: consulAddress,
	}
	sd.ConsulAPI, err = consulapi.NewClient(config)
	if err != nil {
		return sd, err
	}
	sd.ConsulClient = consul.NewClient(sd.ConsulAPI)
	return sd, nil
}

// ConsulInstance creates kit consul instancer which is used to find specific service
// For each service a new instance is required
func (sd *ServiceDiscovery) ConsulInstance(serviceName string, tags []string, passingOnly bool) *consul.Instancer {
	instancer := consul.NewInstancer(sd.ConsulClient, sd.Logger, serviceName, tags, passingOnly)
	return instancer
}
