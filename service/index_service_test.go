package service

import (
	"elasticsearch-version-control/testing_utils"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/stretchr/testify/suite"
	"testing"
)

type IndexServiceIntegrationTestSuite struct {
	suite.Suite
	indexService IndexService
	testUtil     testing_utils.TestUtil
	esClient     *elasticsearch.Client
}

func TestIndexServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IndexServiceIntegrationTestSuite))
}

const INDEX_NAME = "test_index_v10"
const TARGET_INDEX = "test_index_v11"

var (
	INDEX_MAPPING = `{
						"dynamic": "strict",
						"properties": {
						  "description": {
							"type": "text"
						  },
						  "b_id": {
							"type": "keyword"
						  }
						}
					  }`
	INDEX_CONFIG = fmt.Sprintf(`{
					  "settings": {
						"index": {
						  "number_of_shards": 3,
						  "number_of_replicas": 2
						}
					  },
					  "mappings": %v
					}`, INDEX_MAPPING)
)

func (suite *IndexServiceIntegrationTestSuite) SetupSuite() {
	suite.esClient, _ = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	suite.indexService = NewIndexService(suite.esClient)
	suite.testUtil = testing_utils.TestUtil{ElasticClient: suite.esClient}
}

func (suite *IndexServiceIntegrationTestSuite) SetupTest() {
}

func (suite *IndexServiceIntegrationTestSuite) TearDownTest() {
	suite.indexService.DeleteIndex(INDEX_NAME)
	suite.indexService.DeleteIndex(TARGET_INDEX)
}

func (suite IndexServiceIntegrationTestSuite) TestShouldCreateIndexWithCorrectMappings() {
	err := suite.indexService.CreateIndex(INDEX_NAME, INDEX_CONFIG)

	suite.Nil(err)
	mapping, mappingError := suite.testUtil.GetMapping(INDEX_NAME)

	suite.Nil(mappingError)
	isIndexPresent, indexPresentError := suite.testUtil.IsIndexPresent(INDEX_NAME)
	suite.Nil(indexPresentError)
	suite.True(isIndexPresent)
	expectedMapping := fmt.Sprintf(`{
							"%v": {
								"mappings": %v
							}
						}`, INDEX_NAME, INDEX_MAPPING)
	suite.JSONEq(expectedMapping, mapping)
}

func (suite IndexServiceIntegrationTestSuite) TestShouldReturnErrorOnCreateIndexIfItsAlreadyPresent() {
	err := suite.indexService.CreateIndex(INDEX_NAME, INDEX_CONFIG)
	suite.Nil(err)

	indexError := suite.indexService.CreateIndex(INDEX_NAME, INDEX_CONFIG)

	suite.Equal(errors.New("could not create index"), indexError)
}

func (suite IndexServiceIntegrationTestSuite) TestShouldReturnCreationErrorWhenIndexConfigurationIsNotValid() {
	err := suite.indexService.CreateIndex(INDEX_NAME, `{"invalid_config": true}`)
	suite.NotNil(err)
}

func (suite IndexServiceIntegrationTestSuite) TestShouldDeleteIndex() {
	suite.indexService.CreateIndex(INDEX_NAME, INDEX_CONFIG)
	indexPresent, indexPresentError := suite.testUtil.IsIndexPresent(INDEX_NAME)
	suite.Nil(indexPresentError)
	suite.True(indexPresent)

	deletionError := suite.indexService.DeleteIndex(INDEX_NAME)

	suite.Nil(deletionError)
	isIndexPresent, isIndexPresentError := suite.testUtil.IsIndexPresent(INDEX_NAME)
	suite.Nil(isIndexPresentError)
	suite.False(isIndexPresent)
}

func (suite IndexServiceIntegrationTestSuite) TestShouldReturnNumberOfDocumentsForIndex() {
	createIndexErr := suite.indexService.CreateIndex(INDEX_NAME, INDEX_CONFIG)
	suite.Nil(createIndexErr)

	createErr1 := suite.testUtil.CreateDocument(INDEX_NAME, "1", `{"description": "hello", "b_id": "10"}`)
	createErr2 := suite.testUtil.CreateDocument(INDEX_NAME, "2", `{"description": "hello", "b_id": "20"}`)
	createErr3 := suite.testUtil.CreateDocument(INDEX_NAME, "3", `{"description": "hello", "b_id": "30"}`)
	suite.Nil(createErr1)
	suite.Nil(createErr2)
	suite.Nil(createErr3)

	count, documentsCountErr := suite.indexService.GetDocumentsCount(INDEX_NAME)

	suite.Nil(documentsCountErr)
	suite.Equal(3, count)
}

func (suite IndexServiceIntegrationTestSuite) TestShouldReturnErrorWhenIndexIsNotPresent() {
	count, documentsCountErr := suite.indexService.GetDocumentsCount(INDEX_NAME)

	suite.Equal(errors.New("Cannot find documents for index test_index_v10"), documentsCountErr)
	suite.Equal(0, count)
}
