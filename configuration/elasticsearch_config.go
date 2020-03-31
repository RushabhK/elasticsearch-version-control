package configuration

import (
	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearchConfiguration interface {
	GetESClientConfig() elasticsearch.Config
}

type elasticSearchConfiguration struct {
	elasticSearch ElasticSearch
}

func (config elasticSearchConfiguration) GetESClientConfig() elasticsearch.Config {
	return elasticsearch.Config{Addresses: config.elasticSearch.ClientParams.Addresses}
}

func NewElasticSearchConfiguration(elasticSearch ElasticSearch) ElasticSearchConfiguration {
	return elasticSearchConfiguration{elasticSearch: elasticSearch}
}
