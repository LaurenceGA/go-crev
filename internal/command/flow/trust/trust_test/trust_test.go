package trust_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LaurenceGA/go-crev/internal/command/flow/trust"
	"github.com/LaurenceGA/go-crev/internal/command/flow/trust/mock"
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
	"github.com/LaurenceGA/go-crev/internal/github"
	"github.com/LaurenceGA/go-crev/internal/id"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockConfigResponse struct {
	config *config.Configuration
	err    error
}

func (m *mockConfigResponse) getMock(controller *gomock.Controller) *mock.MockConfigReader {
	mck := mock.NewMockConfigReader(controller)
	mck.EXPECT().Load().Return(m.config, m.err)

	return mck
}

type mockGetUser struct {
	expectedUsername string
	usr              *github.User
	err              error
}

type mockGithubResponse struct {
	mockGetUser *mockGetUser
	mockGetRepo *struct {
		expectedOwner, expectedRepo string
		repo                        *github.Repository
		err                         error
	}
}

func (m *mockGithubResponse) getMock(controller *gomock.Controller) *mock.MockGithub {
	mck := mock.NewMockGithub(controller)

	if m.mockGetUser != nil {
		mck.EXPECT().
			GetUser(gomock.Any(), m.mockGetUser.expectedUsername).
			Return(m.mockGetUser.usr, m.mockGetUser.err)
	}

	if m.mockGetRepo != nil {
		mck.EXPECT().
			GetRepository(gomock.Any(), m.mockGetRepo.expectedOwner, m.mockGetRepo.expectedRepo).
			Return(m.mockGetRepo.repo, m.mockGetRepo.err)
	}

	return mck
}

func TestCannotReadConfig(t *testing.T) {
	const (
		testStore = "/my/store"
	)

	for _, testCase := range []struct {
		name               string
		usernameInput      string
		mockConfigResponse mockConfigResponse
		mockGithubResponse mockGithubResponse
		expectError        bool
	}{
		{
			name: "Cannot read config",
			mockConfigResponse: mockConfigResponse{
				err: errors.New("can't read config"),
			},
			expectError: true,
		},
		{
			name: "No store set",
			mockConfigResponse: mockConfigResponse{
				config: &config.Configuration{
					CurrentStore: "",
					CurrentID:    &id.ID{},
				},
			},
			expectError: true,
		},
		{
			name: "No ID set",
			mockConfigResponse: mockConfigResponse{
				config: &config.Configuration{
					CurrentStore: testStore,
					CurrentID:    nil,
				},
			},
			expectError: true,
		},
		{
			name:          "Error getting user",
			usernameInput: "user",
			mockConfigResponse: mockConfigResponse{
				config: &config.Configuration{
					CurrentStore: testStore,
					CurrentID:    &id.ID{},
				},
			},
			mockGithubResponse: mockGithubResponse{
				mockGetUser: &mockGetUser{
					expectedUsername: "user",
					err:              errors.New("failed to talk to Github"),
				},
			},
			expectError: true,
		},
	} {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)

			trustCreator := trust.NewTrustCreator(
				&io.IO{},
				testCase.mockConfigResponse.getMock(controller),
				testCase.mockGithubResponse.getMock(controller),
			)

			err := trustCreator.CreateTrust(context.Background(), testCase.usernameInput, trust.CreatorOptions{})

			if testCase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			controller.Finish()
		})
	}
}
