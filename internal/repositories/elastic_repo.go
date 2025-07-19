package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"sync"

	"UploadDocument-Saas/config"
	"UploadDocument-Saas/internal/models"
)

// IndexDocument concurrently indexes a document in Elasticsearch
func IndexDocument(ctx context.Context, doc models.Document, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	client := config.GetElasticClient()
	body, _ := json.Marshal(doc)
	res, err := client.Index(
		"documents",
		bytes.NewReader(body),
		client.Index.WithContext(ctx),
	)
	if err != nil {
		errCh <- err
		return
	}
	defer res.Body.Close()
	log.Println("Indexed document:", doc.ID)
}

// SearchDocuments concurrently searches documents in Elasticsearch
func SearchDocuments(ctx context.Context, query map[string]interface{}, wg *sync.WaitGroup, docsCh chan<- models.Document, errCh chan<- error) {
	defer wg.Done()
	client := config.GetElasticClient()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		errCh <- err
		return
	}
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex("documents"),
		client.Search.WithBody(&buf),
	)
	if err != nil {
		errCh <- err
		return
	}
	defer res.Body.Close()
	var r struct {
		Hits struct {
			Hits []struct {
				Source models.Document `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		errCh <- err
		return
	}
	for _, hit := range r.Hits.Hits {
		docsCh <- hit.Source
	}
}
