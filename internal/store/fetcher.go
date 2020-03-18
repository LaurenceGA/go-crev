package store

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/meta"
)

type GitCloner interface {
	Clone(ctx context.Context, url string, location string) (*git.Repository, error)
}

type FileDirs interface {
	Cache() string
}

func NewFetcher(cloner GitCloner, fileDirs FileDirs) *Fetcher {
	return &Fetcher{
		gitCloner: cloner,
		fileDirs:  fileDirs,
	}
}

type Fetcher struct {
	gitCloner GitCloner
	fileDirs  FileDirs
}

const (
	cacheDir    = "store"
	cacheGitDir = "git"
)

// Fetch will download a store from a URL to the cache.
func (f *Fetcher) Fetch(ctx context.Context, fetchURL string) error {
	//TODO get repo name
	_, err := f.gitCloner.Clone(ctx, fetchURL, f.cacheDir()+"test")
	if err != nil {
		return fmt.Errorf("cloning git repo: %w", err)
	}

	return nil
}

func (f *Fetcher) cacheDir() string {
	return filepath.Join(f.fileDirs.Cache(), meta.AppName, cacheDir, cacheGitDir)
}
