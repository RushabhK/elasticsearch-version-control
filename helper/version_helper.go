package helper

import (
	"elasticsearch-version-control/model"
)

type VersionHelper interface {
	GetNextIndexVersions(alias string, currentVersion int) ([]model.IndexDetails, error)
}

type versionHelper struct {
	pathToMigrations string
}

func NewVersionHelper(pathToMigrations string) VersionHelper {
	return versionHelper{pathToMigrations: pathToMigrations}
}

func (versionHelper versionHelper) GetNextIndexVersions(alias string, currentVersion int) ([]model.IndexDetails, error) {
	panic("Error")
}
