package trust_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LaurenceGA/go-crev/internal/command/flow/trust"
	"github.com/LaurenceGA/go-crev/internal/command/io"
	"github.com/LaurenceGA/go-crev/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCannotReadConfig(t *testing.T) {
	controller := gomock.NewController(t)

	mockConfigReader := mocks.NewMockConfigReader(controller)

	mockConfigReader.EXPECT().Load().Return(nil, errors.New("can't read config"))

	trustCreator := trust.NewTrustCreator(&io.IO{}, mockConfigReader)

	err := trustCreator.CreateTrust(context.Background(), "", trust.CreatorOptions{})
	assert.Error(t, err)

	controller.Finish()
}
