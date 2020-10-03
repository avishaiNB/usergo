package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ByIDRequest ...
type ByIDRequest struct {
	ID int `json:"id"`
}

// GetUserResponse ...
type GetUserResponse struct {
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

// NewGetUserResponse ...
func NewGetUserResponse(user User) GetUserResponse {
	return GetUserResponse{ID: user.ID, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName}
}

// encoding the response into json
// e.g. GetUsetByIDResponse --> json
func encodeReponseToJSON(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// decoding request into object (acting as anti corruption layer)
// e.g. url --> GetUserByIDRequest
func decodeGetUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	req := NewByIDRequest(id)
	fmt.Println(req)
	return req, nil
}
