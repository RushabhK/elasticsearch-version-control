package service

import (
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"net/http"
	"strings"
)

type AliasService interface {
	SetAlias(alias string, index string) error
	GetIndexVersion(alias string) (int, error)
}

type aliasService struct {
	esClient *elasticsearch.Client
}

func NewAliasService(esClient *elasticsearch.Client) AliasService {
	return aliasService{
		esClient: esClient,
	}
}

func (aliasService aliasService) SetAlias(alias string, index string) error {
	requestBody := fmt.Sprintf(`{
	"actions" : [
			{ "remove" : { "index" : "*",  "alias" : "%v" } },
			{ "add" : { "index" : "%v", "alias" : "%v" } }
		]
	}`, alias, index, alias)
	reader := strings.NewReader(requestBody)

	response, aliasError := aliasService.esClient.Indices.UpdateAliases(reader)
	if aliasError != nil {
		return aliasError
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("could not set alias")
	}
	return nil
}

func (aliasService) GetIndexVersion(alias string) (int, error) {
	panic("implement me")
}
