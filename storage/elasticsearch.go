package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/elastic/go-elasticsearch"
)

type ElasticSearchStorage struct {
	index  string
	client *elasticsearch.Client
}

func (s *ElasticSearchStorage) List() {
	type product struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Price       float32   `json:"price"`
		Available   time.Time `json:"available"`
	}
	var products []product

	response, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(s.index),
	)
	if err != nil {
		return
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return
	}

	for _, hit := range result.Hits.Hits {
		products = append(products, hit.Source)
	}

	err = response.Body.Close()

	return
}

func (s *ElasticSearchStorage) Write() {}

func NewElasticSearchStorage(index string, client *elasticsearch.Client) ElasticSearchStorage {
	return ElasticSearchStorage{
		index:  index,
		client: client,
	}
}
