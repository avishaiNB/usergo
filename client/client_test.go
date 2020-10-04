package client

import (
	"context"
	"testing"
)

func TestClientIntegration(t *testing.T) {
	client, _ := NewServiceClient()
	ctx := context.Background()
	id := 1

	response, err := client.GetUserByID(ctx, id)

	if err != nil {
		t.Error(err)
	}

	if response.User.ID != id {
		t.Error("ID not equal to 1")
	}
}
