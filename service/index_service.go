package service

import (
	"errors"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"net/http"
	"strings"
)

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

func (indexService indexService) CreateIndex(indexName string, configuration string) error {
	response, creationError := indexService.esClient.Indices.Create(indexName, func(request *esapi.IndicesCreateRequest) {
		request.Index = indexName
		request.Body = strings.NewReader(configuration)
	})

	if creationError != nil {
		return creationError
	}
	if response.StatusCode != http.StatusOK {
		return errors.New("could not create index")
	}
	return nil
}

func (indexService) ReIndex(sourceIndex string, targetIndex string, script string) (int, error) {
	panic("implement me")
}

func (indexService indexService) DeleteIndex(indexName string) error {
	response, deleteError := indexService.esClient.Indices.Delete([]string{indexName})
	if deleteError != nil {
		return deleteError
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return errors.New("could not delete index: " + indexName)
	}
	return nil
}

func (indexService) GetDocumentsCount(indexName string) (int, error) {
	panic("implement me")
}
