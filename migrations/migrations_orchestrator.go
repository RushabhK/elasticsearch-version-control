package migrations

import (
	"elasticsearch-version-control/helper"
	"elasticsearch-version-control/service"
)

type MigrationsOrchestrator interface {
	Migrate(alias string) error
}

type migrationsOrchestrator struct {
	indexService  service.IndexService
	aliasService  service.AliasService
	versionHelper helper.VersionHelper
}

func NewMigrationsOrchestrator(indexService service.IndexService, aliasService service.AliasService, versionHelper helper.VersionHelper) MigrationsOrchestrator {
	return migrationsOrchestrator{indexService: indexService, aliasService: aliasService, versionHelper: versionHelper}
}

func (orchestrator migrationsOrchestrator) Migrate(alias string) error {
	panic("Error")
}
