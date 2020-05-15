package git

import (
	"context"
	"errors"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"gopkg.in/src-d/go-git.v4"
)

type Repository struct {
	repo *git.Repository
}

// NewClient construct a new client.
func NewClient(commandIO *io.IO) *Client {
	return &Client{
		commandIO: commandIO,
	}
}

// Client wraps the go-git API.
type Client struct {
	commandIO *io.IO
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const ErrRepositoryAlreadyExists Error = "Repository already exists"

func (g *Client) Clone(ctx context.Context, url string, location string) (*Repository, error) {
	repo, err := git.PlainCloneContext(ctx, location, false, &git.CloneOptions{
		URL:      url,
		Progress: g.commandIO.Out(),
	})
	if err != nil {
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return nil, ErrRepositoryAlreadyExists
		}

		return nil, err
	}

	return &Repository{
		repo: repo,
	}, nil
}
