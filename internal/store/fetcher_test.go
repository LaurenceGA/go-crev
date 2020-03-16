package store

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FetcherSuite struct {
	suite.Suite
}

func TestFetcher(t *testing.T) {
	suite.Run(t, &FetcherSuite{})
}
