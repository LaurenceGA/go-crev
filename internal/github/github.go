package github

import (
	"context"

	"github.com/google/go-github/github"
)

func NewClient() *Client {
	return &Client{
		client: github.NewClient(nil),
	}
}

type Client struct {
	client *github.Client
}

type User struct {
	ID    int64
	Login string
}

func (c *Client) GetUser(ctx context.Context, login string) (*User, error) {
	usr, _, err := c.client.Users.Get(ctx, login)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:    usr.GetID(),
		Login: usr.GetLogin(),
	}, nil
}
