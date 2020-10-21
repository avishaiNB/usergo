package servicediscovery

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/dnssrv"
	consulapi "github.com/hashicorp/consul/api"
	tleerrors "github.com/thelotter-enterprise/usergo/core/errors"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

const (
	// DefaultTTL is 30 seconds
	DefaultTTL time.Duration = time.Second * 30
)

// ServiceDiscovery provide service discovery capabilities for consul and for DNS (k8s)
type ServiceDiscovery struct {
	ConsulAPI      *consulapi.Client
	ConsulClient   *consul.Client
	Logger         log.Logger
	ConsulIntances map[string]*consul.Instancer
	DNSIntances    map[string]*dnssrv.Instancer
}

// NewServiceDiscovery creates a new instance of the service directory
func NewServiceDiscovery(logger log.Logger) ServiceDiscovery {
	sd := ServiceDiscovery{
		Logger:         logger,
		ConsulIntances: map[string]*consul.Instancer{},
		DNSIntances:    map[string]*dnssrv.Instancer{},
	}
	return sd
}

// WithConsul builds consul client and add it to out service discovery
func (sd *ServiceDiscovery) WithConsul(consulAddress string) error {
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
// It will cache the instances
func (sd *ServiceDiscovery) ConsulInstance(serviceName string, tags []string, onlyHealthy bool) (*consul.Instancer, error) {
	key := utils.NewKeys().Build("consul", serviceName, tags...)

	var instancer *consul.Instancer = sd.ConsulIntances[key]
	if instancer != nil {
		return instancer, nil
	}

	if *sd.ConsulClient == nil {
		err := tleerrors.New("call WithConsul first")
		return instancer, err
	}

	instancer = consul.NewInstancer(*sd.ConsulClient, sd.Logger, serviceName, tags, onlyHealthy)
	sd.ConsulIntances[key] = instancer

	return instancer, nil
}

// DNSInstance will return DNS instancer which will be used to lookup a DNS service
// It will cache the instances
func (sd *ServiceDiscovery) DNSInstance(serviceName string) *dnssrv.Instancer {
	key := utils.NewKeys().Build("dns", serviceName)

	var instancer *dnssrv.Instancer = sd.DNSIntances[key]
	if instancer != nil {
		return instancer
	}

	instancer = dnssrv.NewInstancer(serviceName, DefaultTTL, sd.Logger)
	sd.DNSIntances[key] = instancer

	return instancer
}
