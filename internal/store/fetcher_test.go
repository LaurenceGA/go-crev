package store

import (
	"context"
	"errors"
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

	fetcher := NewFetcher(mockGitCloner)

	err := fetcher.Fetch(context.Background(), "")

	s.Error(err)
}

func (s *FetcherSuite) TestCloneSuccess() {
	mockGitCloner := mocks.NewMockGitCloner(s.controller)
	mockGitCloner.EXPECT().
		Clone(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)

	fetcher := NewFetcher(mockGitCloner)

	err := fetcher.Fetch(context.Background(), "")

	s.NoError(err)
}
