package service

import "github.com/elastic/go-elasticsearch/v7"

type IndexService interface {
	CreateIndex(indexName string, configuration string) error
	ReIndex(sourceIndex string, targetIndex string, script string) (int, error)
	DeleteIndex(indexName string) error
	GetDocumentsCount(indexName string) (int, error)
}

type indexService struct {
	esClient *elasticsearch.Client
}

func NewIndexService(esClient *elasticsearch.Client) IndexService {
	return indexService{
		esClient: esClient,
	}
}

func (indexService) CreateIndex(indexName string, configuration string) error {
	panic("implement me")
}

func (indexService) ReIndex(sourceIndex string, targetIndex string, script string) (int, error) {
	panic("implement me")
}

func (indexService) DeleteIndex(indexName string) error {
	panic("implement me")
}

func (indexService) GetDocumentsCount(indexName string) (int, error) {
	panic("implement me")
}
