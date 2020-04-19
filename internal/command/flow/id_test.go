package flow

import (
	"context"
	"errors"
	"testing"

	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/git"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/id"
	"github.com/LaurenceGA/go-crev/internal/mocks"
	"github.com/LaurenceGA/go-crev/internal/store"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type getUserMock struct {
	expectedInput string
	usr           *github.User
	err           error
}

type getRepoInput struct {
	owner, repo string
}

type getRepoMock struct {
	expectedInput getRepoInput
	repo          *github.Repository
	err           error
}

type fetchRepoMock struct {
	expectedFetchURL string
	proofStore       *store.ProofStore
	err              error
}

type setIDMock struct {
	expectedSetID *id.ID
	err           error
}

type setStoreMock struct {
	expectedSetStore string
	err              error
}

type idFlowTestCase struct {
	name          string
	inputUsername string
	getUserMock   getUserMock
	getRepoMock   *getRepoMock
	fetchRepoMock *fetchRepoMock
	setIDMock     *setIDMock
	setStoreMock  *setStoreMock
	expectedError bool
}

func TestIDFlow(t *testing.T) {
	for _, testCase := range []idFlowTestCase{
		{
			name:          "simple success",
			inputUsername: "123",
			getUserMock: getUserMock{
				expectedInput: "123",
				usr:           &github.User{ID: 502, Login: "123"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("123"),
				repo:          &github.Repository{CloneURL: "cloneURL", HTMLurl: "htmlURL"},
			},
			fetchRepoMock: &fetchRepoMock{expectedFetchURL: "cloneURL"},
			setIDMock:     &setIDMock{expectedSetID: githubIDWithURL("502", "htmlURL")},
		},
		{
			name:          "username with @",
			inputUsername: "@user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 18, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           github.NotFoundError,
			},
			setIDMock: &setIDMock{expectedSetID: githubID("18")},
		},
		{
			name:          "username with double @",
			inputUsername: "@@user",
			getUserMock: getUserMock{
				expectedInput: "@user",
				usr:           &github.User{ID: 5, Login: "@user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("@user"),
				err:           github.NotFoundError,
			},
			setIDMock: &setIDMock{expectedSetID: githubID("5")},
		},
		{
			name:          "Get user error",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				err:           errors.New("not found"),
			},
			expectedError: true,
		},
		{
			name:          "Fail to set ID",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 5, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           github.NotFoundError,
			},
			setIDMock: &setIDMock{
				expectedSetID: githubID("5"),
				err:           errors.New("failed to update config"),
			},
			expectedError: true,
		},
		{
			name:          "Fail to get repo from github",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 5, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           errors.New("that's a fail"),
			},
			setIDMock:     &setIDMock{expectedSetID: githubID("5")},
			expectedError: false, // Failing to find repo is non-fatal
		},
		{
			name:          "Repo already exists",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 5, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				repo:          &github.Repository{CloneURL: "cloneURL", HTMLurl: "htmlURL"},
			},
			fetchRepoMock: &fetchRepoMock{
				expectedFetchURL: "cloneURL",
				err:              git.ErrRepositoryAlreadyExists,
			},
			setIDMock:     &setIDMock{expectedSetID: githubIDWithURL("5", "htmlURL")},
			expectedError: false, // Failing to find repo is non-fatal
		},
		{
			name:          "Fail to clone repo",
			inputUsername: "user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 5, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				repo:          &github.Repository{CloneURL: "cloneURL", HTMLurl: "htmlURL"},
			},
			fetchRepoMock: &fetchRepoMock{
				expectedFetchURL: "cloneURL",
				err:              errors.New("That's a fail"),
			},
			setIDMock:     &setIDMock{expectedSetID: githubIDWithURL("5", "htmlURL")},
			expectedError: false, // Failing to find repo is non-fatal
		},
	} {
		t.Run(testCase.name, runIDFlowTestCase(testCase))
	}
}

func runIDFlowTestCase(testCase idFlowTestCase) func(*testing.T) {
	return func(t *testing.T) {
		controller := gomock.NewController(t)
		mockConfigManipulator := mocks.NewMockConfigManipulator(controller)
		mockGithub := mocks.NewMockGithub(controller)
		mockRepoFetcher := mocks.NewMockRepoFetcher(controller)

		setGetUserMocks(mockGithub, testCase.getUserMock)
		setGetRepoMocks(mockGithub, testCase.getRepoMock)
		setIDSetMocks(mockConfigManipulator, testCase.setIDMock)
		setStoreSetMocks(mockConfigManipulator, testCase.setStoreMock)
		setFetchRepoMocks(mockRepoFetcher, testCase.fetchRepoMock)

		idSetterFlow := NewIDSetter(&io.IO{}, mockConfigManipulator, mockGithub, mockRepoFetcher)

		err := idSetterFlow.SetFromUsername(context.Background(), testCase.inputUsername)

		if testCase.expectedError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}

		controller.Finish()
	}
}

func setGetUserMocks(mockGithub *mocks.MockGithub, mockData getUserMock) {
	mockGithub.EXPECT().
		GetUser(gomock.Any(), mockData.expectedInput).
		Return(mockData.usr, mockData.err)
}

func setGetRepoMocks(mockGithub *mocks.MockGithub, mockData *getRepoMock) {
	if mockData != nil {
		mockGithub.EXPECT().
			GetRepository(gomock.Any(), mockData.expectedInput.owner, mockData.expectedInput.repo).
			Return(mockData.repo, mockData.err)
	}
}

func setIDSetMocks(mockConfigManipulator *mocks.MockConfigManipulator, mockData *setIDMock) {
	if mockData != nil {
		mockConfigManipulator.EXPECT().
			SetCurrentID(mockData.expectedSetID).
			Return(mockData.err)
	}
}

func setStoreSetMocks(mockConfigManipulator *mocks.MockConfigManipulator, mockData *setStoreMock) {
	if mockData != nil {
		mockConfigManipulator.EXPECT().
			SetCurrentStore(mockData.expectedSetStore).
			Return(mockData.err)
	}
}

func setFetchRepoMocks(mockRepoFetcher *mocks.MockRepoFetcher, mockData *fetchRepoMock) {
	if mockData != nil {
		mockRepoFetcher.EXPECT().
			Fetch(gomock.Any(), mockData.expectedFetchURL).
			Return(mockData.proofStore, mockData.err)
	}
}

func expectedDefaultRepo(owner string) getRepoInput {
	return getRepoInput{
		owner: owner,
		repo:  standardCrevProofRepoName,
	}
}

func githubID(githubID string) *id.ID {
	return &id.ID{
		ID:   githubID,
		Type: id.Github,
	}
}

func githubIDWithURL(githubID, url string) *id.ID {
	return &id.ID{
		ID:   githubID,
		Type: id.Github,
		URL:  url,
	}
}
