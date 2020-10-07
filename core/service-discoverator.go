package core

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
)

// this should be infra

// ServiceDiscoverator ...
type ServiceDiscoverator struct {
	ConsulAPI    *consulapi.Client
	ConsulClient *consul.Client
	Logger       log.Logger
}

// NewServiceDiscovery ...
func NewServiceDiscovery(logger log.Logger) ServiceDiscoverator {
	sd := ServiceDiscoverator{
		Logger: logger,
	}
	return sd
}

// WithConsul builds consul client and add it to out service discovery
func (sd *ServiceDiscoverator) WithConsul(consulAddress string) error {
	var err error
	config := &consulapi.Config{
		Address: consulAddress,
	}

	sd.ConsulAPI, err = consulapi.NewClient(config)

	if err == nil {
		client := consul.NewClient(sd.ConsulAPI)
		sd.ConsulClient = &client
	} else {
		sd.Logger.Log("method", "NewServiceDiscovery", "input", consulAddress, "err", err)
	}

	return err
}

// ConsulInstance creates kit consul instancer which is used to find specific service
// For each service a new instance is required
func (sd *ServiceDiscoverator) ConsulInstance(serviceName string, tags []string, passingOnly bool) (*consul.Instancer, error) {
	var instancer *consul.Instancer
	if *sd.ConsulClient == nil {
		err := NewApplicationError("call WithConsul first", nil)
		return instancer, err
	}

	instancer = consul.NewInstancer(*sd.ConsulClient, sd.Logger, serviceName, tags, passingOnly)
	return instancer, nil
}
