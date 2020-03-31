package model

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type IndexDetailsTestSuite struct {
	suite.Suite
}

func TestIndexDetailsTestSuite(t *testing.T) {
	suite.Run(t, new(IndexDetailsTestSuite))
}

func (suite IndexDetailsTestSuite) TestShouldGetNameOfIndexFromIndexDetails() {
	indexDetails := IndexDetails{Alias: "index", Version: 15}

	suite.Equal("index_v15", indexDetails.GetName())
}
