package fetcher

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/LaurenceGA/go-crev/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type FetcherSuite struct {
	suite.Suite

	controller *gomock.Controller
}

func (s *FetcherSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
}

func (s *FetcherSuite) TearDownTest() {
	s.controller.Finish()
}

func TestFetcher(t *testing.T) {
	suite.Run(t, &FetcherSuite{})
}

func (s *FetcherSuite) TestFailToClone() {
	mockGitCloner := mocks.NewMockGitCloner(s.controller)
	mockGitCloner.EXPECT().
		Clone(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("failed to clone"))

	mockFileDirs := mocks.NewMockFileDirs(s.controller)
	mockFileDirs.EXPECT().Data().Return("data", nil)

	fetcher := NewFetcher(mockGitCloner, mockFileDirs)

	_, err := fetcher.Fetch(context.Background(), "")

	s.Error(err)
}

func (s *FetcherSuite) TestFailToFindCloneDir() {
	mockGitCloner := mocks.NewMockGitCloner(s.controller)
	mockFileDirs := mocks.NewMockFileDirs(s.controller)
	mockFileDirs.EXPECT().Data().Return("", errors.New("no filesystem"))

	fetcher := NewFetcher(mockGitCloner, mockFileDirs)

	_, err := fetcher.Fetch(context.Background(), "")

	s.Error(err)
}

func (s *FetcherSuite) TestCloneSuccess() {
	mockGitCloner := mocks.NewMockGitCloner(s.controller)
	mockGitCloner.EXPECT().
		Clone(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)

	mockFileDirs := mocks.NewMockFileDirs(s.controller)
	mockFileDirs.EXPECT().Data().Return("data", nil)

	fetcher := NewFetcher(mockGitCloner, mockFileDirs)

	proofStore, err := fetcher.Fetch(context.Background(), "github.com/path")

	s.NoError(err)
	s.Equal(filepath.Join("data", "store", "git", "github.com", "path"), proofStore.Dir)
}

func (s *FetcherSuite) TestURLToPath() {
	for _, tt := range []struct {
		name, url, expectedPath string
	}{
		{
			name:         "raw URL",
			url:          "plain",
			expectedPath: "plain",
		},
		{
			name:         "a slash",
			url:          "parent/child",
			expectedPath: filepath.Join("parent", "child"),
		},
		{
			name:         "with protocol, only hostname",
			url:          "https://example",
			expectedPath: "example",
		},
		{
			name:         "with protocol, hostname and path",
			url:          "https://example.com/path",
			expectedPath: filepath.Join("example.com", "path"),
		},
		{
			name:         "longer path",
			url:          "https://example.com/one/two/three",
			expectedPath: filepath.Join("example.com", "one", "two", "three"),
		},
		{
			name:         "git protocol",
			url:          "git:git.example.com/octocat/Hello-World",
			expectedPath: filepath.Join("git.example.com", "octocat", "Hello-World"),
		},
		{
			name:         "SSH protocol",
			url:          "git@github.com:octocat/Hello-World.git",
			expectedPath: filepath.Join("github.com", "octocat", "Hello-World"),
		},
		{
			name:         "longer SSH protocol",
			url:          "ssh://server/project.git",
			expectedPath: filepath.Join("server", "project"),
		},
		{
			name:         "local file",
			url:          "/srv/git/project.git",
			expectedPath: filepath.Join("srv", "git", "project"),
		},
		{
			name:         "longer local file",
			url:          "file:///srv/git/project.git",
			expectedPath: filepath.Join("srv", "git", "project"),
		},
	} {
		tt := tt
		s.Run(tt.name, func() {
			path, err := pathFromRepoURL(tt.url)
			s.NoError(err) // Effectively never returns an error, so we won't test it

			s.Equal(tt.expectedPath, path)
		})
	}
}
