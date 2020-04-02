package helper

import (
	"elasticsearch-version-control/model"
)

type VersionDetailsProvider interface {
	GetNextIndexVersions(alias string, currentVersion int) ([]model.IndexDetails, error)
}

type versionDetailsProvider struct {
	pathToMigrations string
}

func NewVersionDetailsProvider(pathToMigrations string) VersionDetailsProvider {
	return versionDetailsProvider{pathToMigrations: pathToMigrations}
}

func (versionDetailsProvider versionDetailsProvider) GetNextIndexVersions(alias string, currentVersion int) ([]model.IndexDetails, error) {
	panic("Implement me")
}
