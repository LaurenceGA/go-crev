package fetcher

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/store"
	giturls "github.com/whilp/git-urls"
)

type GitCloner interface {
	Clone(ctx context.Context, url string, location string) (*git.Repository, error)
}

type FileDirs interface {
	Data() (string, error)
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
	storeDataDir = "store"
	storeGitDir  = "git"
)

// Fetch will download a store from a URL to the cache.
func (f *Fetcher) Fetch(ctx context.Context, fetchURL string) (*store.ProofStore, error) {
	repoPath, err := pathFromRepoURL(fetchURL)
	if err != nil {
		return nil, fmt.Errorf("converting URL '%s' into path: %w", fetchURL, err)
	}

	dataDir, err := f.dataDir()
	if err != nil {
		return nil, fmt.Errorf("finding data directory: %w", err)
	}

	cloneDir := filepath.Join(dataDir, repoPath)
	fmt.Printf("Cloning into %s\n", cloneDir)

	if _, err := f.gitCloner.Clone(ctx, fetchURL, cloneDir); err != nil {
		return nil, fmt.Errorf("cloning git repo: %w", err)
	}

	return &store.ProofStore{
		Dir: cloneDir,
	}, nil
}

const gitProtocolURLExtension = ".git"

func pathFromRepoURL(repoURL string) (string, error) {
	repoURL = strings.TrimSuffix(repoURL, gitProtocolURLExtension)

	u, err := giturls.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("parsing repo URL %s: %w", repoURL, err)
	}

	rawHostnamePath := filepath.Join(u.Hostname(), u.EscapedPath())
	if rawHostnamePath == "" {
		// Rootless paths (E.G. git:example.com/repo) are considered opaque, so we just use that
		rawHostnamePath = u.Opaque
	}

	hostnamePath := strings.TrimPrefix(rawHostnamePath, "/") // Must be relative

	return filepath.FromSlash(hostnamePath), nil
}

func (f *Fetcher) dataDir() (string, error) {
	dataDir, err := f.fileDirs.Data()
	if err != nil {
		return "", err
	}

	return filepath.Join(dataDir, storeDataDir, storeGitDir), nil
}
