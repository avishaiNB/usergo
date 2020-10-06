package client

import "github.com/thelotter-enterprise/usergo/shared"

// Endpoints ...
type Endpoints struct {
	EP map[string]*shared.ProxyEndpoint
}

// NewEndpoints ..
func NewEndpoints() Endpoints {
	return Endpoints{
		EP: make(map[string]*shared.ProxyEndpoint),
	}
}

// Add will add a new endpoint
func (endpoints *Endpoints) Add(name string, endpoint shared.ProxyEndpoint) {
	endpoints.EP[name] = &endpoint
}

// Get will return the endpoint based on the name
func (endpoints *Endpoints) Get(name string) shared.ProxyEndpoint {
	return *endpoints.EP[name]
}
