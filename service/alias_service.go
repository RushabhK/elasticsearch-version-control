package service

import (
	"github.com/elastic/go-elasticsearch/v7"
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

func (aliasService) SetAlias(alias string, index string) error {
	panic("implement me")
}

func (aliasService) GetIndexVersion(alias string) (int, error) {
	panic("implement me")
}
