package user

import "context"

// Service ...
type Service interface {
	GetUserByID(ctx context.Context, userID int) (User, error)
}
