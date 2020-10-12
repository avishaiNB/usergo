package core

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// EncodeRequestToJSON is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func EncodeRequestToJSON(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// EncodeReponseToJSON will encoding the response into json
// e.g. GetUsetByIDResponse --> json
func EncodeReponseToJSON(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// ParsePath ...
// TODO: implement the parse: /v1/user/{id} --> /v1/user/123
func ParsePath(path string, data ...interface{}) {

}

// DecodeRequestFromJSON ....
func DecodeRequestFromJSON(ctx context.Context, r *http.Request) (interface{}, error) {
	var req Request
	d, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(d, &req)
	//err := json.NewDecoder(r.Body).Decode(&req)

	return req, err
}
