package testing_utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"net/http"
	"strings"
	"time"
)

type TestUtil struct {
	ElasticClient *elasticsearch.Client
}

func (testUtil TestUtil) GetMapping(index string) (string, error) {
	response, err := testUtil.ElasticClient.Indices.GetMapping(func(request *esapi.IndicesGetMappingRequest) {
		request.Index = []string{index}
	})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	return buf.String(), nil
}

func (testUtil TestUtil) IsIndexPresent(index string) (bool, error) {
	response, err := testUtil.ElasticClient.Indices.Get([]string{index})
	if err != nil {
		return false, err
	}
	return response.StatusCode != http.StatusNotFound, nil
}

func (testUtil TestUtil) CreateDocument(index string, id string, body string) error {
	response, err := testUtil.ElasticClient.Create(index, id, strings.NewReader(body), func(request *esapi.CreateRequest) {
		request.Timeout = 5 * time.Second
	})
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return errors.New("Could not create document")
	}
	testUtil.ElasticClient.Indices.Refresh()

	return nil
}

func (testUtil TestUtil) GetIndexSetByAlias(alias string) (string, error) {
	response, err := testUtil.ElasticClient.Indices.GetAlias(func(request *esapi.IndicesGetAliasRequest) {
		request.Name = []string{alias}
	})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	var aliasResponse map[string]interface{}
	json.Unmarshal(buf.Bytes(), &aliasResponse)
	var index string
	for k := range aliasResponse {
		index = k
	}
	if index == "" {
		return index, errors.New("No index associated with alias")
	}
	return index, nil
}

func (testUtil TestUtil) GetDocument(index string, id string) (string, error) {
	response, err := testUtil.ElasticClient.Get(index, id)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", errors.New("Did not find document")
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	document := esDoc{}
	json.Unmarshal(buf.Bytes(), &document)
	content, _ := json.Marshal(document.Source)
	return string(content), nil
}

type esDoc struct {
	Source map[string]interface{} `json:"_source"`
}
