package client

import (
	"github.com/thelotter-enterprise/usergo/core"
)

// Endpoints ...
type Endpoints struct {
	EP map[string]*core.ProxyEndpoint
}

// NewEndpoints ..
func NewEndpoints() Endpoints {
	return Endpoints{
		EP: make(map[string]*core.ProxyEndpoint),
	}
}

// Add will add a new endpoint
func (endpoints *Endpoints) Add(name string, endpoint core.ProxyEndpoint) {
	endpoints.EP[name] = &endpoint
}

// Get will return the endpoint based on the name
func (endpoints *Endpoints) Get(name string) core.ProxyEndpoint {
	return *endpoints.EP[name]
}
