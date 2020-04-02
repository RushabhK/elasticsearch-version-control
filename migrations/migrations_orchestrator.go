package migrations

import (
	migrations_error "elasticsearch-version-control/error"
	"elasticsearch-version-control/helper"
	"elasticsearch-version-control/service"
	"errors"
	"fmt"
	"strconv"
)

type MigrationsOrchestrator interface {
	Migrate(alias string) error
}

type migrationsOrchestrator struct {
	indexService  service.IndexService
	aliasService  service.AliasService
	versionHelper helper.VersionDetailsProvider
}

func NewMigrationsOrchestrator(indexService service.IndexService, aliasService service.AliasService, versionHelper helper.VersionDetailsProvider) MigrationsOrchestrator {
	return migrationsOrchestrator{indexService: indexService, aliasService: aliasService, versionHelper: versionHelper}
}

func (orchestrator migrationsOrchestrator) Migrate(alias string) error {
	currentIndexVersion, aliasError := orchestrator.aliasService.GetIndexVersion(alias)
	currentIndexName := alias + "_v" + strconv.Itoa(currentIndexVersion)
	var sourceIndexDocCount int
	if aliasError == migrations_error.AliasNotFoundError {
		currentIndexVersion = 0
		sourceIndexDocCount = 0
	} else if aliasError != nil {
		return aliasError
	} else {
		var countErr error
		sourceIndexDocCount, countErr = orchestrator.indexService.GetDocumentsCount(currentIndexName)
		if countErr != nil {
			return countErr
		}
	}
	indexDetails, nextVersionsError := orchestrator.versionHelper.GetNextIndexVersions(alias, currentIndexVersion)
	if nextVersionsError != nil {
		return nextVersionsError
	}
	for _, nextIndexVersionDetail := range indexDetails {
		deleteError := orchestrator.indexService.DeleteIndex(nextIndexVersionDetail.GetName())
		if deleteError != nil {
			return deleteError
		}
		createError := orchestrator.indexService.CreateIndex(nextIndexVersionDetail.GetName(), nextIndexVersionDetail.IndexConfig)
		if createError != nil {
			return createError
		}
		if nextIndexVersionDetail.GetName() != alias+"_v1" {
			createdDocCount, reIndexError := orchestrator.indexService.ReIndex(currentIndexName, nextIndexVersionDetail.GetName(), nextIndexVersionDetail.Script)
			if reIndexError != nil {
				return reIndexError
			}
			if createdDocCount != sourceIndexDocCount {
				errorMsg := fmt.Sprintf("All documents not migrated from %s to %s", currentIndexName, nextIndexVersionDetail.GetName())
				return errors.New(errorMsg)
			}
		}
		currentIndexName = nextIndexVersionDetail.GetName()
	}
	return orchestrator.aliasService.SetAlias(alias, currentIndexName)
}
