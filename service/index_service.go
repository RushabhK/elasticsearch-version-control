package service

import "github.com/elastic/go-elasticsearch/v7"

type IndexService interface{}

type indexService struct {
	esClient *elasticsearch.Client
}

func NewIndexService(esClient *elasticsearch.Client) IndexService {
	return indexService{
		esClient: esClient,
	}
}
