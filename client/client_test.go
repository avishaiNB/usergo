package client

import (
	"context"
	"testing"
)

func TestClientIntegration(t *testing.T) {
	client, _ := NewServiceClient()
	ctx := context.Background()
	id := 1

	response := client.GetUserByID(ctx, id)

	if response.Result != nil {
		t.Fail()
	}
}
