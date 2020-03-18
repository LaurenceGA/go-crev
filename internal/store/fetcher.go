package store

import (
	"context"
	"fmt"

	"github.com/LaurenceGA/go-crev/internal/git"
)

type GitCloner interface {
	Clone(ctx context.Context, url string, location string) (*git.Repository, error)
}

func NewFetcher(cloner GitCloner) *Fetcher {
	return &Fetcher{
		gitCloner: cloner,
	}
}

type Fetcher struct {
	gitCloner GitCloner
}

// Fetch will download a store from a URL to the cache.
func (f *Fetcher) Fetch(ctx context.Context, fetchURL string) error {
	_, err := f.gitCloner.Clone(ctx, fetchURL, "./repo/test")
	if err != nil {
		return fmt.Errorf("cloning git repo: %w", err)
	}

	return nil
}
