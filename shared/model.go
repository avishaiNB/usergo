package shared

import (
	"context"

	"github.com/thelotter-enterprise/usergo/core"
)

// User ...
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ByIDRequestData ...
type ByIDRequestData struct {
	ID int `json:"id"`
}

// ByIDResponseData ...
type ByIDResponseData struct {
	User User
}

// NewByIDRequest will create a Request with ByIDRequestData
func NewByIDRequest(ctx context.Context, id int) core.Request {
	data := ByIDRequestData{
		ID: id,
	}
	req := core.Request{}.Wrap(ctx, data)
	return req
}

// NewUserResponse ...
func NewUserResponse(user User) ByIDResponseData {
	return ByIDResponseData{User: user}
}
