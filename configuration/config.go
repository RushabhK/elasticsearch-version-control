package configuration

import (
	"encoding/json"
	"os"
	"time"
)

type Config interface {
	LoadConfig(fileName string) (*ConfigData, error)
}

type configLoader struct {
}

func NewConfigLoader() Config {
	return configLoader{}
}

type ConfigData struct {
	ElasticSearch ElasticSearch `json:"elasticsearch"`
}

type ElasticSearch struct {
	ClientParams            ClientParams  `json:"client_params"`
	ReindexTimeoutInMinutes time.Duration `json:"reindex_timeout_in_minutes"`
}

type ClientParams struct {
	Addresses []string `json:"addresses"`
}

func (configLoader configLoader) LoadConfig(fileName string) (*ConfigData, error) {
	configuration := ConfigData{}
	file, err := os.Open(fileName)
	if err != nil {
		return &configuration, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	return &configuration, err
}
