package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		testname string
		username string
	}{
		{
			testname: "getUser alice (user exists, if migration ...02 applied)",
			username: "alice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			t.Parallel()

			_, err := getUser(context.Background(), tt.username)
			assert.NoError(t, err)
		})
	}
}
