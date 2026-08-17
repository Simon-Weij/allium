package users

import (
	"context"
	"fmt"

	"github.com/Simon-Weij/allium/generated/sqlc"
)


func GetUserByUsername(ctx context.Context, username string, queries *sqlc.Queries) (*sqlc.User, error) {
	user, err := queries.GetUser(ctx, username)
	if err != nil {
		return &sqlc.User{}, fmt.Errorf("error recieved when trying to query user from Database: %w", err)
	}
	return &user, nil
}

func GetUsers(ctx context.Context, queries *sqlc.Queries) (*[]sqlc.User, error) {
	users, err := queries.GetUsers(ctx)
	if err != nil {
		return &[]sqlc.User{}, fmt.Errorf("error recieved when trying to get all users from: %w", err)
	}
	return &users, nil
}