package github

import (
	"context"
	"net/http"

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

type Error string

const NotFoundError Error = "Not Found"

func (e Error) Error() string {
	return string(e)
}

type User struct {
	ID    int64
	Login string
}

func (c *Client) GetUser(ctx context.Context, login string) (*User, error) {
	usr, resp, err := c.client.Users.Get(ctx, login)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, NotFoundError
		}

		return nil, err
	}

	return &User{
		ID:    usr.GetID(),
		Login: usr.GetLogin(),
	}, nil
}

type Repository struct {
	CloneURL string
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*Repository, error) {
	repository, resp, err := c.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, NotFoundError
		}

		return nil, err
	}

	return &Repository{
		CloneURL: repository.GetCloneURL(),
	}, nil
}
