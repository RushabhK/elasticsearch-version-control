package testing_utils

import (
	"bytes"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"net/http"
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
