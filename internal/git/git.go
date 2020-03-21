package git

import (
	"context"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

type Repository struct {
	repo *git.Repository
}

// NewClient construct a new client
func NewClient() *Client {
	return &Client{}
}

// Client wraps the go-git API
type Client struct {
}

func (g *Client) Clone(ctx context.Context, url string, location string) (*Repository, error) {
	repo, err := git.PlainCloneContext(ctx, location, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,	//TODO change to cmd output
	})

	if err != nil {
		return nil, err
	}

	return &Repository{
		repo: repo,
	}, nil
}
