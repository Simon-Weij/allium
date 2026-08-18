package users

import (
	"context"
	"fmt"

	"github.com/Simon-Weij/allium/generated/sqlc"
)

type UserClient struct {
	Queries *sqlc.Queries
}

func NewUserClient(q *sqlc.Queries) *UserClient {
	return &UserClient{
		Queries: q,
	}
}

// Similarly to iTunes, make a struct that satisfies the UserManagementClient Interface, use that for the mocks instead
// Then Figure out how we're going to deal with the folder array (json_array() sql, to text in query)

func (u *UserClient) GetUserByUsername(ctx context.Context, username string) (*sqlc.User, error) {
	user, err := u.Queries.GetUser(ctx, username)
	if err != nil {
		return &sqlc.User{}, fmt.Errorf("error recieved when trying to query user from Database: %w", err)
	}
	return &user, nil
}

func (u *UserClient) GetUsers(ctx context.Context) (*[]sqlc.User, error) {
	users, err := u.Queries.GetUsers(ctx)
	if err != nil {
		return &[]sqlc.User{}, fmt.Errorf("error recieved when trying to get all users from: %w", err)
	}
	return &users, nil
}