package svc

import (
	"context"

	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
	"github.com/thelotter-enterprise/usergo/shared"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
)

// ServiceMiddleware used to chain behaviors on the UserService using middleware pattern
type ServiceMiddleware func(Service) Service

// Service API
type Service interface {
	GetUserByID(ctx context.Context, userID int) (shared.User, error)
	ConsumeLoginCommand(ctx context.Context, userID int) error
}

type service struct {
	repo   Repository
	tracer tletracer.Tracer
	logger tlelogger.Manager
}

// NewService creates a new instance of service
// service is where we define all the business logic.
func NewService(logger tlelogger.Manager, tracer tletracer.Tracer, repo Repository) Service {
	return &service{
		repo:   repo,
		tracer: tracer,
		logger: logger,
	}
}

// GetUserByID will execute business logic for getting user information by id
func (s *service) GetUserByID(ctx context.Context, userID int) (shared.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	return user, err
}

func (s *service) ConsumeLoginCommand(ctx context.Context, userID int) error {
	return nil
}
