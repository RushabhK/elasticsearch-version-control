package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type IndexService interface {
	CreateIndex(indexName string, configuration string) error
	ReIndex(sourceIndex string, targetIndex string, script string) (int, error)
	DeleteIndex(indexName string) error
	GetDocumentsCount(indexName string) (int, error)
}

type indexService struct {
	esClient                *elasticsearch.Client
	reindexTimeoutInMinutes time.Duration
}

func NewIndexService(esClient *elasticsearch.Client, reindexTimeoutInMinutes time.Duration) IndexService {
	return indexService{
		esClient:                esClient,
		reindexTimeoutInMinutes: reindexTimeoutInMinutes,
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

func (indexService indexService) ReIndex(sourceIndex string, targetIndex string, script string) (int, error) {
	reader := indexService.getReindexRequestBody(sourceIndex, targetIndex, script)
	waitForCompletion := true
	reIndexRequest := esapi.ReindexRequest{
		Body:              reader,
		Pretty:            true,
		WaitForCompletion: &waitForCompletion,
		Timeout:           indexService.reindexTimeoutInMinutes * time.Minute,
	}

	response, reIndexError := reIndexRequest.Do(nil, indexService.esClient)

	if reIndexError != nil {
		logrus.Error("Error while reindexing: ", reIndexError.Error())
		return 0, reIndexError
	} else if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("could not reindex from %v to %v", sourceIndex, targetIndex)
		logrus.Error(errorMessage)
		return 0, errors.New(errorMessage)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	_, refreshError := indexService.esClient.Indices.Refresh()
	if refreshError != nil {
		return 0, refreshError
	}
	responseString := buf.String()
	logrus.Info("Reindexing source: ", sourceIndex, ", targetIndex: ", targetIndex, ", Response: ", responseString)
	var responseObj map[string]interface{}
	json.Unmarshal(buf.Bytes(), &responseObj)
	createdDocCount := int(responseObj["created"].(float64))
	return createdDocCount, nil
}

func (indexService indexService) getReindexRequestBody(sourceIndex, targetIndex, script string) *strings.Reader {
	var requestBody string
	if script == "" {
		requestBody = fmt.Sprintf(`{
									"source": {
										"index": "%v"
									},
									"dest": {
										"index": "%v"
									}
								}`, sourceIndex, targetIndex)
	} else {
		requestBody = fmt.Sprintf(`{
							"source": {
								"index": "%v"
							},
							"dest": {
								"index": "%v"
							},
							"script": {
								"lang": "painless",
								"source": "%v"
							}
						}`, sourceIndex, targetIndex, script)
	}
	return strings.NewReader(requestBody)
}
