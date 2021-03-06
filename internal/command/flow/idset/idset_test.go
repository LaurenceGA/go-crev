package idset

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

type fetchStoreMock struct {
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

func TestIDFlow(t *testing.T) {
	const (
		testUser            = "user"
		testUserIDNum int64 = 502
		testUserIDStr       = "502"
		testCloneURL        = "cloneURL"
	)

	for name, testCase := range map[string]struct {
		inputUsername  string
		getUserMock    getUserMock
		getRepoMock    *getRepoMock
		fetchStoreMock *fetchStoreMock
		setIDMock      *setIDMock
		setStoreMock   *setStoreMock
		expectedError  bool
	}{
		"simple success": {
			inputUsername: "user123",
			getUserMock: getUserMock{
				expectedInput: "user123",
				usr:           &github.User{ID: 5, Login: "user123"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user123"),
				repo:          &github.Repository{CloneURL: "someCloningURL", HTMLurl: "htmlURL"},
			},
			fetchStoreMock: &fetchStoreMock{
				expectedFetchURL: "someCloningURL",
				proofStore:       &store.ProofStore{Dir: "store/path"},
			},
			setStoreMock: &setStoreMock{
				expectedSetStore: "store/path",
			},
			setIDMock: &setIDMock{expectedSetID: githubIDWithURL("5", "user123", "htmlURL")},
		},
		"username with @": {
			inputUsername: "@user",
			getUserMock: getUserMock{
				expectedInput: "user",
				usr:           &github.User{ID: 18, Login: "user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("user"),
				err:           github.NotFoundError,
			},
			setIDMock: &setIDMock{expectedSetID: githubID("18", "user")},
		},
		"username with double @": {
			inputUsername: "@@user",
			getUserMock: getUserMock{
				expectedInput: "@user",
				usr:           &github.User{ID: testUserIDNum, Login: "@user"},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo("@user"),
				err:           github.NotFoundError,
			},
			setIDMock: &setIDMock{expectedSetID: githubID(testUserIDStr, "@user")},
		},
		"Get user error": {
			inputUsername: testUser,
			getUserMock: getUserMock{
				expectedInput: testUser,
				err:           errors.New("not found"),
			},
			expectedError: true,
		},
		"Fail to set ID": {
			inputUsername: testUser,
			getUserMock: getUserMock{
				expectedInput: testUser,
				usr:           &github.User{ID: testUserIDNum, Login: testUser},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo(testUser),
				err:           github.NotFoundError,
			},
			setIDMock: &setIDMock{
				expectedSetID: githubID(testUserIDStr, testUser),
				err:           errors.New("failed to update config"),
			},
			expectedError: true,
		},
		"Fail to get repo from github": {
			inputUsername: testUser,
			getUserMock: getUserMock{
				expectedInput: testUser,
				usr:           &github.User{ID: testUserIDNum, Login: testUser},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo(testUser),
				err:           errors.New("that's a fail"),
			},
			setIDMock:     &setIDMock{expectedSetID: githubID(testUserIDStr, testUser)},
			expectedError: false, // Failing to find repo is non-fatal
		},
		"Repo already exists": {
			inputUsername: testUser,
			getUserMock: getUserMock{
				expectedInput: testUser,
				usr:           &github.User{ID: testUserIDNum, Login: testUser},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo(testUser),
				repo:          &github.Repository{CloneURL: testCloneURL, HTMLurl: "htmlURL"},
			},
			fetchStoreMock: &fetchStoreMock{
				expectedFetchURL: testCloneURL,
				proofStore:       &store.ProofStore{Dir: "store/path"},
				err:              git.ErrRepositoryAlreadyExists,
			},
			setStoreMock: &setStoreMock{
				expectedSetStore: "store/path",
			},
			setIDMock:     &setIDMock{expectedSetID: githubIDWithURL(testUserIDStr, testUser, "htmlURL")},
			expectedError: false, // Failing to find repo is non-fatal
		},
		"Fail to clone repo": {
			inputUsername: testUser,
			getUserMock: getUserMock{
				expectedInput: testUser,
				usr:           &github.User{ID: testUserIDNum, Login: testUser},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo(testUser),
				repo:          &github.Repository{CloneURL: testCloneURL, HTMLurl: "htmlURL"},
			},
			fetchStoreMock: &fetchStoreMock{
				expectedFetchURL: testCloneURL,
				err:              errors.New("That's a fail"),
			},
			setIDMock:     &setIDMock{expectedSetID: githubIDWithURL(testUserIDStr, testUser, "htmlURL")},
			expectedError: false, // Failing to find repo is non-fatal
		},
		"Fail to set Store": {
			inputUsername: testUser,
			getUserMock: getUserMock{
				expectedInput: testUser,
				usr:           &github.User{ID: testUserIDNum, Login: testUser},
			},
			getRepoMock: &getRepoMock{
				expectedInput: expectedDefaultRepo(testUser),
				repo:          &github.Repository{CloneURL: testCloneURL, HTMLurl: "repoURL"},
			},
			fetchStoreMock: &fetchStoreMock{
				expectedFetchURL: testCloneURL,
				proofStore:       &store.ProofStore{Dir: "store/path"},
				err:              git.ErrRepositoryAlreadyExists,
			},
			setStoreMock: &setStoreMock{
				expectedSetStore: "store/path",
				err:              errors.New("it failed"),
			},
			setIDMock:     &setIDMock{expectedSetID: githubIDWithURL(testUserIDStr, testUser, "repoURL")},
			expectedError: false, // Error is non-fatal
		},
	} {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			controller := gomock.NewController(t)
			mockConfigManipulator := mocks.NewMockConfigManipulator(controller)
			mockGithub := mocks.NewMockGithub(controller)
			mockRepoFetcher := mocks.NewMockRepoFetcher(controller)

			setGetUserMocks(mockGithub, testCase.getUserMock)
			setGetRepoMocks(mockGithub, testCase.getRepoMock)
			setIDSetMocks(mockConfigManipulator, testCase.setIDMock)
			setStoreSetMocks(mockConfigManipulator, testCase.setStoreMock)
			setFetchStoreMocks(mockRepoFetcher, testCase.fetchStoreMock)

			idSetterFlow := NewIDSetter(&io.IO{}, mockConfigManipulator, mockGithub, mockRepoFetcher)

			err := idSetterFlow.SetFromUsername(context.Background(), testCase.inputUsername)

			if testCase.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
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

func setFetchStoreMocks(mockStoreFetcher *mocks.MockRepoFetcher, mockData *fetchStoreMock) {
	if mockData != nil {
		mockStoreFetcher.EXPECT().
			Fetch(gomock.Any(), mockData.expectedFetchURL).
			Return(mockData.proofStore, mockData.err)
	}
}

func expectedDefaultRepo(owner string) getRepoInput {
	return getRepoInput{
		owner: owner,
		repo:  store.StandardCrevProofRepoName,
	}
}

func githubID(githubID, alias string) *id.ID {
	return &id.ID{
		ID:    githubID,
		Type:  id.Github,
		Alias: alias,
	}
}

func githubIDWithURL(githubID, alias, url string) *id.ID {
	return &id.ID{
		ID:    githubID,
		Type:  id.Github,
		URL:   url,
		Alias: alias,
	}
}
