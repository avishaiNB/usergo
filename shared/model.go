package shared

import "context"

// HTTPResponse is a base response
type HTTPResponse struct {
	Error         error
	Result        interface{}
	CircuitOpened bool
	Context       context.Context
	StatusCode    int
	// CorrelactionID
	// timeout
}

// User ...
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// ByIDRequest ...
type ByIDRequest struct {
	ID int `json:"id"`
}

// ByIDResponse ...
type ByIDResponse struct {
	User User
}

// NewByIDRequest ...
func NewByIDRequest(id int) ByIDRequest {
	return ByIDRequest{
		ID: id,
	}
}

// NewUserResponse ...
func NewUserResponse(user User) ByIDResponse {
	return ByIDResponse{User: user}
}
