package configuration

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ElasticSearchConfigTestSuite struct {
	suite.Suite
}

func TestElasticSearchConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ElasticSearchConfigTestSuite))
}

func (suite ElasticSearchConfigTestSuite) TestShouldGetESClientConfigs() {
	configuration := NewElasticSearchConfiguration(ElasticSearch{
		ClientParams{Addresses: []string{"localhost:9200"}},
		15,
	})

	suite.Equal([]string{"localhost:9200"}, configuration.GetESClientConfig().Addresses)
}
