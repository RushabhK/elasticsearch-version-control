package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
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

func (indexService indexService) GetDocumentsCount(indexName string) (int, error) {
	response, err := indexService.esClient.Count(func(request *esapi.CountRequest) {
		request.Index = []string{indexName}
	})
	if err != nil {
		return 0, nil
	}
	if response.StatusCode != http.StatusOK {
		errorMsg := "Cannot find documents for index " + indexName
		logrus.Error(errorMsg)
		return 0, errors.New(errorMsg)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	var responseObj map[string]interface{}
	json.Unmarshal(buf.Bytes(), &responseObj)
	count := responseObj["count"]
	return int(count.(float64)), nil
}
