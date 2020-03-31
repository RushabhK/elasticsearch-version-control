package service

import (
	"github.com/elastic/go-elasticsearch/v7"
)

type AliasService interface{}

type aliasService struct {
	esClient *elasticsearch.Client
}

func NewAliasService(esClient *elasticsearch.Client) AliasService {
	return aliasService{
		esClient: esClient,
	}
}
