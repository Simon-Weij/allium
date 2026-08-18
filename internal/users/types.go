package users

import (
	"context"
	"github.com/Simon-Weij/allium/generated/sqlc"
)

//go:generate mockgen -source=types.go -destination=../../generated/mocks/users_mock.go -package=mocks
type UserManagementClient interface {
	GetUserByUsername(ctx context.Context, username string) (*sqlc.User, error)
	GetUsers(ctx context.Context) (*[]sqlc.User, error)
}