package client_test

import (
	"context"
	"testing"

	"github.com/thelotter-enterprise/usergo/client"
)

func TestClientIntegration(t *testing.T) {
	c := client.NewServiceClient()
	ctx := context.Background()
	id := 1

	response := c.GetUserByID(ctx, id)

	if response.Result == nil {
		t.Fail()
	}
}
