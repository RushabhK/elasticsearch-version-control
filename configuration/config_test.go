package configuration

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ConfigTestSuite struct {
	suite.Suite
	config Config
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) SetupTest() {
	suite.config = NewConfigLoader()
}

func (suite ConfigTestSuite) TestShouldLoadConfigFromFile() {
	configData, err := suite.config.LoadConfig("./config.json")
	suite.Nil(err)
	suite.Equal(time.Duration(15), configData.ElasticSearch.ReindexTimeoutInMinutes)
}

func (suite ConfigTestSuite) TestShouldThrowErrorWhenFileNotPresent() {
	_, err := suite.config.LoadConfig("../resources/randomFile.json")
	suite.NotNil(err)
}
