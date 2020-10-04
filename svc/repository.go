package user

import (
	"context"

	om "github.com/thelotter-enterprise/usergo/shared"
)

// Repository ..
type Repository interface {
	GetUserByID(ctx context.Context, userID int) (om.User, error)
}

type repo struct {
	// database for example
}

// NewRepository ...
func NewRepository() Repository {
	return &repo{}
}

// GetUserByID ...
func (r repo) GetUserByID(ctx context.Context, userID int) (om.User, error) {

	user := om.User{
		ID:        userID,
		Email:     "guyk@net-bet.net",
		FirstName: "guy",
		LastName:  "kolbis",
	}

	return user, nil
}
