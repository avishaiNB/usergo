package user

import "context"

type service struct {
	repo Repository
}

// NewService ...
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetUserByID(ctx context.Context, userID int) (User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	return user, err
}
