package store

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/meta"
	giturls "github.com/whilp/git-urls"
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
	repoPath, err := pathFromRepoURL(fetchURL)
	if err != nil {
		return fmt.Errorf("cloning git repo: %w", err)
	}

	if _, err := f.gitCloner.Clone(ctx, fetchURL, filepath.Join(f.cacheDir(), repoPath)); err != nil {
		return fmt.Errorf("cloning git repo: %w", err)
	}

	return nil
}

func pathFromRepoURL(repoURL string) (string, error) {
	u, err := giturls.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("parsing repo URL %s: %w", repoURL, err)
	}

	return filepath.FromSlash(filepath.Join(u.Hostname(), u.EscapedPath())), nil
}

func (f *Fetcher) cacheDir() string {
	return filepath.Join(f.fileDirs.Cache(), meta.AppName, cacheDir, cacheGitDir)
}
