package store

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/LaurenceGA/go-crev/internal/git"
	giturls "github.com/whilp/git-urls"
)

type GitCloner interface {
	Clone(ctx context.Context, url string, location string) (*git.Repository, error)
}

type FileDirs interface {
	Cache() (string, error)
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
	storeCacheDir = "store"
	cacheGitDir   = "git"
)

// Fetch will download a store from a URL to the cache.
func (f *Fetcher) Fetch(ctx context.Context, fetchURL string) error {
	repoPath, err := pathFromRepoURL(fetchURL)
	if err != nil {
		return fmt.Errorf("converting URL '%s' into path: %w", fetchURL, err)
	}

	cacheDir, err := f.cacheDir()
	if err != nil {
		return fmt.Errorf("finding cache directory: %w", err)
	}

	cloneDir := filepath.Join(cacheDir, repoPath)
	fmt.Printf("Cloning into %s\n", cloneDir)

	if _, err := f.gitCloner.Clone(ctx, fetchURL, cloneDir); err != nil {
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

func (f *Fetcher) cacheDir() (string, error) {
	cacheDir, err := f.fileDirs.Cache()
	if err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, storeCacheDir, cacheGitDir), nil
}
