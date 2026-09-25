package users

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Simon-Weij/allium/generated/sqlc"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

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

func (u *UserClient) GetUserByUsername(ctx context.Context, username string) (*sqlc.GetUserRow, error) {
	user, err := u.Queries.GetUser(ctx, username)
	if err != nil {
		return &sqlc.GetUserRow{}, fmt.Errorf("error recieved when trying to query user from Database: %w", err)
	}
	return &user, nil
}

func (u *UserClient) GetUsers(ctx context.Context) (*[]sqlc.GetUsersRow, error) {
	users, err := u.Queries.GetUsers(ctx)
	if err != nil {
		return &[]sqlc.GetUsersRow{}, fmt.Errorf("error recieved when trying to get all users from: %w", err)
	}
	return &users, nil
}

func (u *UserClient) CreateUser(ctx context.Context, reqValues url.Values) (error) {

	var params sqlc.CreateUserParams
	err := decoder.Decode(&params, reqValues)
	if err != nil {
		return fmt.Errorf("error returned while trying to createUser %w", err)
	}
	err = u.Queries.CreateUser(ctx, params) 
	if err != nil {
		return fmt.Errorf("error returned while trying to createUser %w", err)
	}
	return nil
}

func (u *UserClient) GetPassword(ctx context.Context, username string, q sqlc.Queries) (string, error) {
	pass, err := u.Queries.GetPassword(ctx, username)
	if err != nil {
		return "", err
	}

	return pass, nil
}