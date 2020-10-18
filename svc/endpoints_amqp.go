package svc

import (
	"github.com/thelotter-enterprise/usergo/core"
)

// UserAMQPEndpoints ...
type UserAMQPEndpoints struct {
	Service       Service
	Log           core.Log
	Tracer        core.Tracer
	AMQPEndpoints *core.AMQPEndpoints
}

// NewUserAMQPEndpoints ...
func NewUserAMQPEndpoints(log core.Log, tracer core.Tracer, service Service) *UserAMQPEndpoints {
	userEndpoints := UserAMQPEndpoints{
		Log:     log,
		Tracer:  tracer,
		Service: service,
	}

	userEndpoints.AMQPEndpoints = userEndpoints.makeEndpoints()

	return &userEndpoints
}

func (a UserAMQPEndpoints) makeEndpoints() *core.AMQPEndpoints {
	endpoints := core.NewAMQPEndpoints()
	return &endpoints
}
