package config

import (
	"log"
	"os"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	esClient *elasticsearch.Client
	esOnce   sync.Once
)

// GetElasticClient returns a singleton Elasticsearch client
func GetElasticClient() *elasticsearch.Client {
	esOnce.Do(func() {
		addr := os.Getenv("ELASTIC_URL")
		if addr == "" {
			addr = "http://localhost:9200"
		}
		cfg := elasticsearch.Config{
			Addresses: []string{addr},
		}
		client, err := elasticsearch.NewClient(cfg)
		if err != nil {
			log.Fatalf("Failed to connect to Elasticsearch: %v", err)
		}
		esClient = client
		log.Println("Connected to Elasticsearch")
	})
	return esClient
}
