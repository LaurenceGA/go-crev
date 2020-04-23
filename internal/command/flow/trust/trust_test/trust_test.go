package trust_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LaurenceGA/go-crev/internal/command/flow/trust"
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/config"
	"github.com/LaurenceGA/go-crev/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockConfigResponse struct {
	config *config.Configuration
	err    error
}

func (m *mockConfigResponse) applyToMock(mock *mocks.MockConfigReader) {
	mock.EXPECT().Load().Return(m.config, m.err)
}

func (m *mockConfigResponse) getConfigReader(controller *gomock.Controller) *mocks.MockConfigReader {
	mock := mocks.NewMockConfigReader(controller)
	m.applyToMock(mock)

	return mock
}

func TestCannotReadConfig(t *testing.T) {
	const (
		testStore = "/my/store"
	)

	for _, testCase := range []struct {
		name               string
		mockConfigResponse mockConfigResponse
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
	} {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)

			trustCreator := trust.NewTrustCreator(
				&io.IO{},
				testCase.mockConfigResponse.getConfigReader(controller),
			)

			err := trustCreator.CreateTrust(context.Background(), "", trust.CreatorOptions{})

			if testCase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			controller.Finish()
		})
	}
}
