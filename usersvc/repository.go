package user

import "context"

// Repository ..
type Repository interface {
	GetUserByID(ctx context.Context, userID int) (User, error)
}

type repo struct {
	// database for example
}

// NewRepository ...
func NewRepository() Repository {
	return &repo{}
}

// GetUserByID ...
func (r repo) GetUserByID(ctx context.Context, userID int) (User, error) {

	user := User{
		ID:        userID,
		Email:     "guyk@net-bet.net",
		FirstName: "guy",
		LastName:  "kolbis",
	}

	return user, nil
}
