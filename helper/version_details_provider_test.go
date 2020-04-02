package helper

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type VersionDetailsProviderTestSuite struct {
	suite.Suite
	mockCtrl               *gomock.Controller
	versionDetailsProvider VersionDetailsProvider
}

func TestVersionDetailsProviderTestSuite(t *testing.T) {
	suite.Run(t, new(VersionDetailsProviderTestSuite))
}

func (suite *VersionDetailsProviderTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.versionDetailsProvider = NewVersionDetailsProvider("./test_resources/")
}

func (suite VersionDetailsProviderTestSuite) TestShouldReturnErrorWhenAliasDirectoryIsNotPresent() {
	aliasName := "alias_not_present"
	indexDetails, err := suite.versionDetailsProvider.GetNextIndexVersions(aliasName, 1)

	suite.NotNil(err)
	suite.Nil(indexDetails)
}

func (suite VersionDetailsProviderTestSuite) TestShouldReturnVersion2And3ForNextIndexVersions() {
	aliasName := "test_alias_with_scripts"

	indexDetails, err := suite.versionDetailsProvider.GetNextIndexVersions(aliasName, 1)

	suite.Nil(err)
	suite.Equal(2, len(indexDetails))
	suite.JSONEq(`{"version": "v2"}`, indexDetails[0].IndexConfig)
	suite.JSONEq(`{"version": "v3"}`, indexDetails[1].IndexConfig)
	suite.Equal("script-version2", indexDetails[0].Script)
	suite.Equal("script-version3", indexDetails[1].Script)
	suite.Equal(aliasName, indexDetails[0].Alias)
	suite.Equal(aliasName, indexDetails[1].Alias)
	suite.Equal(2, indexDetails[0].Version)
	suite.Equal(3, indexDetails[1].Version)
}

func (suite VersionDetailsProviderTestSuite) TestShouldReturnEmptyScriptsWhenScriptFileIsNotPresentForVersion() {
	aliasName := "test_alias_without_scripts"

	indexDetails, err := suite.versionDetailsProvider.GetNextIndexVersions(aliasName, 1)

	suite.Nil(err)
	suite.Equal(2, len(indexDetails))
	suite.JSONEq(`{"version": "v2"}`, indexDetails[0].IndexConfig)
	suite.JSONEq(`{"version": "v3"}`, indexDetails[1].IndexConfig)
	suite.Empty(indexDetails[0].Script)
	suite.Equal("script-version3", indexDetails[1].Script)
	suite.Equal(aliasName, indexDetails[0].Alias)
	suite.Equal(aliasName, indexDetails[1].Alias)
	suite.Equal(2, indexDetails[0].Version)
	suite.Equal(3, indexDetails[1].Version)
}

func (suite VersionDetailsProviderTestSuite) TestShouldReturnErrorWhenMappingIsNotAValidJson() {
	aliasName := "test_alias_with_invalid_json"

	indexDetails, err := suite.versionDetailsProvider.GetNextIndexVersions(aliasName, 1)

	suite.Equal(errors.New("./test_resources/test_alias_with_invalid_json/index_configs/v2.json file is not a valid json"), err)
	suite.Nil(indexDetails)
}
