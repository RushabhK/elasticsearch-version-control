package migrations

import (
	migrations_error "elasticsearch-version-control/error"
	"elasticsearch-version-control/mocks"
	"elasticsearch-version-control/model"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MigrationsOrchestratorTestSuite struct {
	suite.Suite
	mockCtrl               *gomock.Controller
	indexService           *mocks.MockIndexService
	aliasService           *mocks.MockAliasService
	versionDetailsProvider *mocks.MockVersionDetailsProvider
	migrationsOrchestrator MigrationsOrchestrator
}

func TestMigrationsOrchestratorTestSuite(t *testing.T) {
	suite.Run(t, new(MigrationsOrchestratorTestSuite))
}

func (suite *MigrationsOrchestratorTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.indexService = mocks.NewMockIndexService(suite.mockCtrl)
	suite.aliasService = mocks.NewMockAliasService(suite.mockCtrl)
	suite.versionDetailsProvider = mocks.NewMockVersionDetailsProvider(suite.mockCtrl)
	suite.migrationsOrchestrator = NewMigrationsOrchestrator(suite.indexService, suite.aliasService, suite.versionDetailsProvider)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldCreateNewIndicesAndSetAliasToLatestVersion() {
	alias := "index_name"
	currentVersion := 1
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(currentVersion, nil)
	suite.indexService.EXPECT().GetDocumentsCount("index_name_v1").Return(3, nil)
	indexDetailsV2 := model.IndexDetails{Version: 2, Script: "version-2-script", IndexConfig: "index-config-2", Alias: alias}
	indexDetailsV3 := model.IndexDetails{Version: 3, Script: "version-3-script", IndexConfig: "index-config-3", Alias: alias}
	indexDetails := []model.IndexDetails{
		indexDetailsV2,
		indexDetailsV3,
	}
	suite.versionDetailsProvider.EXPECT().GetNextIndexVersions(alias, currentVersion).Return(indexDetails, nil)
	suite.indexService.EXPECT().DeleteIndex("index_name_v2").Return(nil)
	suite.indexService.EXPECT().DeleteIndex("index_name_v3").Return(nil)
	suite.indexService.EXPECT().CreateIndex("index_name_v2", "index-config-2").Return(nil)
	suite.indexService.EXPECT().CreateIndex("index_name_v3", "index-config-3").Return(nil)
	suite.indexService.EXPECT().ReIndex("index_name_v1", "index_name_v2", "version-2-script").Return(3, nil)
	suite.indexService.EXPECT().ReIndex("index_name_v2", "index_name_v3", "version-3-script").Return(3, nil)
	suite.aliasService.EXPECT().SetAlias(alias, "index_name_v3").Return(nil)

	migrateError := suite.migrationsOrchestrator.Migrate(alias)

	suite.Nil(migrateError)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldReturnErrorWhenGetIndexVersionFails() {
	alias := "index_name"
	aliasError := errors.New("something went wrong")
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(-1, aliasError)

	migrateError := suite.migrationsOrchestrator.Migrate(alias)

	suite.Equal(aliasError, migrateError)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldCreateIndexAndAliasWhenAliasIsNotPresent() {
	alias := "index_name"
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(-1, migrations_error.AliasNotFoundError)
	indexDetailsV1 := model.IndexDetails{Version: 1, Script: "", IndexConfig: "index-config-1", Alias: alias}
	indexDetailsV2 := model.IndexDetails{Version: 2, Script: "version-2-script", IndexConfig: "index-config-2", Alias: alias}
	indexDetails := []model.IndexDetails{
		indexDetailsV1,
		indexDetailsV2,
	}
	suite.versionDetailsProvider.EXPECT().GetNextIndexVersions(alias, 0).Return(indexDetails, nil)
	suite.indexService.EXPECT().DeleteIndex("index_name_v1").Return(nil)
	suite.indexService.EXPECT().DeleteIndex("index_name_v2").Return(nil)
	suite.indexService.EXPECT().CreateIndex("index_name_v1", "index-config-1").Return(nil)
	suite.indexService.EXPECT().CreateIndex("index_name_v2", "index-config-2").Return(nil)
	suite.indexService.EXPECT().ReIndex("index_name_v1", "index_name_v2", "version-2-script")
	suite.aliasService.EXPECT().SetAlias(alias, "index_name_v2").Return(nil)

	migrateError := suite.migrationsOrchestrator.Migrate(alias)

	suite.Nil(migrateError)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldReturnErrorWhenGetDocumentsCountFails() {
	alias := "index_name"
	currentVersion := 1
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(currentVersion, nil)
	countError := errors.New("count error")
	suite.indexService.EXPECT().GetDocumentsCount("index_name_v1").Return(0, countError)

	migrateError := suite.migrationsOrchestrator.Migrate(alias)

	suite.Equal(countError, migrateError)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldReturnErrorWhenVersionHelperFails() {
	alias := "index_name"
	currentVersion := 1
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(currentVersion, nil)
	suite.indexService.EXPECT().GetDocumentsCount("index_name_v1").Return(1, nil)
	err := errors.New("Could not find folder for migrations")
	suite.versionDetailsProvider.EXPECT().GetNextIndexVersions(alias, currentVersion).Return(nil, err)

	migrateErr := suite.migrationsOrchestrator.Migrate(alias)

	suite.Equal(err, migrateErr)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldReturnErrorWhenDeleteIndexFails() {
	alias := "index_name"
	currentVersion := 1
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(currentVersion, nil)
	suite.indexService.EXPECT().GetDocumentsCount("index_name_v1").Return(1, nil)
	indexDetailsV2 := model.IndexDetails{Version: 2, Script: "version-2-script", IndexConfig: "index-config-2", Alias: alias}
	indexDetails := []model.IndexDetails{
		indexDetailsV2,
	}
	suite.versionDetailsProvider.EXPECT().GetNextIndexVersions(alias, currentVersion).Return(indexDetails, nil)
	deleteError := errors.New("deletion error")
	suite.indexService.EXPECT().DeleteIndex("index_name_v2").Return(deleteError)

	migrateError := suite.migrationsOrchestrator.Migrate(alias)

	suite.Equal(deleteError, migrateError)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldReturnErrorWhenCreateIndexFails() {
	alias := "index_name"
	currentVersion := 1
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(currentVersion, nil)
	suite.indexService.EXPECT().GetDocumentsCount("index_name_v1").Return(1, nil)
	indexDetailsV2 := model.IndexDetails{Version: 2, Script: "version-2-script", IndexConfig: "index-config-2", Alias: alias}
	indexDetails := []model.IndexDetails{
		indexDetailsV2,
	}
	suite.versionDetailsProvider.EXPECT().GetNextIndexVersions(alias, currentVersion).Return(indexDetails, nil)
	suite.indexService.EXPECT().DeleteIndex("index_name_v2").Return(nil)
	createError := errors.New("creation error")
	suite.indexService.EXPECT().CreateIndex("index_name_v2", "index-config-2").Return(createError)

	migrateError := suite.migrationsOrchestrator.Migrate(alias)

	suite.Equal(createError, migrateError)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldReturnErrorWhenReIndexFails() {
	alias := "index_name"
	currentVersion := 1
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(currentVersion, nil)
	suite.indexService.EXPECT().GetDocumentsCount("index_name_v1").Return(1, nil)
	indexDetailsV2 := model.IndexDetails{Version: 2, Script: "version-2-script", IndexConfig: "index-config-2", Alias: alias}
	indexDetails := []model.IndexDetails{
		indexDetailsV2,
	}
	suite.versionDetailsProvider.EXPECT().GetNextIndexVersions(alias, currentVersion).Return(indexDetails, nil)
	suite.indexService.EXPECT().DeleteIndex("index_name_v2").Return(nil)
	suite.indexService.EXPECT().CreateIndex("index_name_v2", "index-config-2").Return(nil)
	reIndexError := errors.New("reIndex error")
	suite.indexService.EXPECT().ReIndex("index_name_v1", "index_name_v2", "version-2-script").Return(0, reIndexError)

	migrateError := suite.migrationsOrchestrator.Migrate(alias)

	suite.Equal(reIndexError, migrateError)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldReturnErrorWhenReindexedCountDoesNotMatchSourceIndexCount() {
	alias := "index_name"
	currentVersion := 1
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(currentVersion, nil)
	sourceIndexCount := 5
	suite.indexService.EXPECT().GetDocumentsCount("index_name_v1").Return(sourceIndexCount, nil)
	indexDetailsV2 := model.IndexDetails{Version: 2, Script: "version-2-script", IndexConfig: "index-config-2", Alias: alias}
	indexDetails := []model.IndexDetails{
		indexDetailsV2,
	}
	suite.versionDetailsProvider.EXPECT().GetNextIndexVersions(alias, currentVersion).Return(indexDetails, nil)
	suite.indexService.EXPECT().DeleteIndex("index_name_v2").Return(nil)
	suite.indexService.EXPECT().CreateIndex("index_name_v2", "index-config-2").Return(nil)
	suite.indexService.EXPECT().ReIndex("index_name_v1", "index_name_v2", "version-2-script").Return(3, nil)

	migrateError := suite.migrationsOrchestrator.Migrate(alias)

	expectedError := errors.New("All documents not migrated from index_name_v1 to index_name_v2")
	suite.Equal(expectedError, migrateError)
}

func (suite MigrationsOrchestratorTestSuite) TestShouldReturnErrorWhenSetAliasFails() {
	alias := "index_name"
	currentVersion := 1
	suite.aliasService.EXPECT().GetIndexVersion(alias).Return(currentVersion, nil)
	suite.indexService.EXPECT().GetDocumentsCount("index_name_v1").Return(1, nil)
	indexDetailsV2 := model.IndexDetails{Version: 2, Script: "version-2-script", IndexConfig: "index-config-2", Alias: alias}
	indexDetails := []model.IndexDetails{
		indexDetailsV2,
	}
	suite.versionDetailsProvider.EXPECT().GetNextIndexVersions(alias, currentVersion).Return(indexDetails, nil)
	suite.indexService.EXPECT().DeleteIndex("index_name_v2").Return(nil)
	suite.indexService.EXPECT().CreateIndex("index_name_v2", "index-config-2").Return(nil)
	suite.indexService.EXPECT().ReIndex("index_name_v1", "index_name_v2", "version-2-script").Return(1, nil)
	setAliasError := errors.New("set alias error")
	suite.aliasService.EXPECT().SetAlias(alias, "index_name_v2").Return(setAliasError)

	migrateError := suite.migrationsOrchestrator.Migrate(alias)

	suite.Equal(setAliasError, migrateError)
}
