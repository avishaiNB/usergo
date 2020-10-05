package svc

import (
	"context"

	"github.com/thelotter-enterprise/usergo/shared"
)

// Service API
type Service interface {
	GetUserByID(ctx context.Context, userID int) (shared.User, error)
}

type service struct {
	repo   Repository
	tracer Tracer
	logger Logger
}

// NewService creates a new instance of service
// service is where we define all the business logic.
func NewService(logger Logger, tracer Tracer, repo Repository) Service {
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
