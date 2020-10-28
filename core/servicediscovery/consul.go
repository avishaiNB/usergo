package servicediscovery

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

// ConsulServiceDiscovery provide service discovery capabilities for consul
type ConsulServiceDiscovery struct {
	ConsulAPI      *consulapi.Client
	ConsulClient   *consul.Client
	Logger         tlelogger.Manager
	ConsulAddress  string
	ConsulIntances map[string]*consul.Instancer
}

// NewConsulServiceDiscovery creates a new instance of the service directory
func NewConsulServiceDiscovery(logger tlelogger.Manager, consulAddress string) ConsulServiceDiscovery {
	sd := ConsulServiceDiscovery{
		Logger:         logger,
		ConsulAddress:  consulAddress,
		ConsulIntances: map[string]*consul.Instancer{},
	}

	return sd
}

// ConsulInstance creates kit consul instancer which is used to find specific service
// For each service a new instance is required
// It will cache the instances
func (sd *ConsulServiceDiscovery) ConsulInstance(ctx context.Context, serviceName string, tags []string, onlyHealthy bool) (*consul.Instancer, error) {
	key := utils.NewKeys().Build("consul", serviceName, tags...)

	var instancer *consul.Instancer = sd.ConsulIntances[key]
	if instancer != nil {
		return instancer, nil
	}

	config := &consulapi.Config{
		Address: sd.ConsulAddress,
	}

	api, err := consulapi.NewClient(config)

	if err == nil {
		client := consul.NewClient(sd.ConsulAPI)
		instancer = consul.NewInstancer(client, sd.Logger.(log.Logger), serviceName, tags, onlyHealthy)

		sd.ConsulAPI = api
		sd.ConsulClient = &client
		sd.ConsulIntances[key] = instancer
	}

	return instancer, err
}
