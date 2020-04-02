package helper

import (
	"elasticsearch-version-control/model"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
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
	pathToAlias := versionDetailsProvider.pathToMigrations + alias
	pathToIndexConfigs := pathToAlias + "/index_configs/"
	pathToScripts := pathToAlias + "/scripts/"
	file, err := os.Open(pathToIndexConfigs)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	defer file.Close()

	indexConfigFiles, err := file.Readdirnames(0)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	var indexDetails []model.IndexDetails

	for version := currentVersion + 1; true; version += 1 {
		fileVersion := "v" + strconv.Itoa(version)
		indexConfigFile := fileVersion + ".json"
		if !isFilePresent(indexConfigFile, indexConfigFiles) {
			break
		}
		indexConfigBytes, _ := ioutil.ReadFile(pathToIndexConfigs + indexConfigFile)
		var jsonObj map[string]interface{}
		if json.Unmarshal(indexConfigBytes, &jsonObj) != nil {
			errorMsg := pathToIndexConfigs + indexConfigFile + " file is not a valid json"
			logrus.Error(errorMsg)
			return nil, errors.New(errorMsg)
		}
		indexConfig := string(indexConfigBytes)
		scriptFile := fileVersion + ".painless"
		scriptBytes, _ := ioutil.ReadFile(pathToScripts + scriptFile)
		script := string(scriptBytes)
		nextIndexDetails := model.IndexDetails{
			Version:     version,
			Alias:       alias,
			IndexConfig: indexConfig,
			Script:      script,
		}
		indexDetails = append(indexDetails, nextIndexDetails)
	}
	return indexDetails, nil
}

func isFilePresent(file string, indexConfigFiles []string) bool {
	for _, indexConfigFile := range indexConfigFiles {
		if file == indexConfigFile {
			return true
		}
	}
	return false
}
