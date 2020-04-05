package service

import (
	"elasticsearch-version-control/testing_utils"
	"errors"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AliasServiceTestSuite struct {
	suite.Suite
	indexService IndexService
	aliasService AliasService
	testUtil     testing_utils.TestUtil
	esClient     *elasticsearch.Client
}

const (
	INDEX = "test_index"
	ALIAS = "test_alias"
)

func (suite *AliasServiceTestSuite) SetupSuite() {
	suite.esClient, _ = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	suite.indexService = NewIndexService(suite.esClient, 1)
	suite.testUtil = testing_utils.TestUtil{ElasticClient: suite.esClient}
}

func TestAliasServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AliasServiceTestSuite))
}

func (suite AliasServiceTestSuite) TestShouldSetAliasToIndex() {
	createIndexErr := suite.indexService.CreateIndex(INDEX, INDEX_CONFIG)
	suite.Nil(createIndexErr)

	aliasError := suite.aliasService.SetAlias(ALIAS, INDEX)

	suite.Nil(aliasError)
	index, err := suite.testUtil.GetIndexSetByAlias(ALIAS)
	suite.Nil(err)
	suite.Equal(INDEX, index)
}

func (suite AliasServiceTestSuite) TestShouldRemoveExistingAliasAndSetAliasToNewIndex() {
	createIndexErr := suite.indexService.CreateIndex(INDEX, INDEX_CONFIG)
	suite.Nil(createIndexErr)
	newIndex := "latest_index_version"
	createIndexErr = suite.indexService.CreateIndex(newIndex, INDEX_CONFIG)
	suite.Nil(createIndexErr)

	aliasError := suite.aliasService.SetAlias(ALIAS, INDEX)
	suite.Nil(aliasError)
	aliasError = suite.aliasService.SetAlias(ALIAS, newIndex)
	suite.Nil(aliasError)
	index, err := suite.testUtil.GetIndexSetByAlias(ALIAS)
	suite.Nil(err)
	suite.Equal(newIndex, index)
	suite.indexService.DeleteIndex(newIndex)
}

func (suite AliasServiceTestSuite) TestShouldReturnErrorWhenAliasIsSetToIndexWhichDoesNotExist() {
	aliasError := suite.aliasService.SetAlias(ALIAS, "index-not-present")
	suite.Equal(errors.New("could not set alias"), aliasError)
}
