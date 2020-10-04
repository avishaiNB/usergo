package om

// User ...
type User struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
}

// ByIDRequest ...
type ByIDRequest struct {
	ID int `json:"id"`
}

// ByIDResponse ...
type ByIDResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// NewByIDRequest ...
func NewByIDRequest(id int) ByIDRequest {
	return ByIDRequest{
		ID: id,
	}
}

// NewUserResponse ...
func NewUserResponse(user User) ByIDResponse {
	return ByIDResponse{ID: user.ID, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName}
}
