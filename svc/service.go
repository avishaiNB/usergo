package svc

import (
	"context"

	"github.com/thelotter-enterprise/usergo/core"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
	"github.com/thelotter-enterprise/usergo/shared"
)

// UserServiceMiddleware used to chain behaviors on the UserService using middleware pattern
type ServiceMiddleware func(Service) Service

// Service API
type Service interface {
	GetUserByID(ctx context.Context, userID int) (shared.User, error)
	ConsumeLoginCommand(ctx context.Context, userID int) error
}

type service struct {
	repo   Repository
	tracer tletracer.Tracer
	log    core.Log
}

// NewService creates a new instance of service
// service is where we define all the business logic.
func NewService(log core.Log, tracer tletracer.Tracer, repo Repository) Service {
	return &service{
		repo:   repo,
		tracer: tracer,
	}
}

// GetUserByID will execute business logic for getting user information by id
func (s *service) GetUserByID(ctx context.Context, userID int) (shared.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	return user, err
}

func (s *service) ConsumeLoginCommand(ctx context.Context, userID int) error {
	s.log.Logger.Log("message", "login command consumed")
	return nil
}
