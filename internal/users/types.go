package users

import (
	"context"
	"net/url"

	"github.com/Simon-Weij/allium/generated/sqlc"
)

//go:generate mockgen -source=types.go -destination=../../generated/mocks/users_mock.go -package=mocks
type UserManagementClient interface {
	GetUserByUsername(ctx context.Context, username string) (*sqlc.GetUserRow, error)
	GetUsers(ctx context.Context) (*[]sqlc.GetUsersRow, error)
	CreateUser(ctx context.Context, reqValues url.Values) (error)
}