package repositories

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"UploadDocument-Saas/config"
	"UploadDocument-Saas/internal/models"
)

var (
	documentCollection *mongo.Collection
	docOnce            sync.Once
)

func getDocumentCollection() *mongo.Collection {
	docOnce.Do(func() {
		client := config.GetMongoClient()
		documentCollection = client.Database("testdb").Collection("documents")
	})
	return documentCollection
}

// InsertDocument inserts a document concurrently
func InsertDocument(ctx context.Context, doc models.Document, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	_, err := getDocumentCollection().InsertOne(ctx, doc)
	if err != nil {
		errCh <- err
		return
	}
	log.Println("Inserted document:", doc.ID)
}

// FindDocuments concurrently finds all documents
func FindDocuments(ctx context.Context, filter bson.M, wg *sync.WaitGroup, docsCh chan<- models.Document, errCh chan<- error) {
	defer wg.Done()
	cur, err := getDocumentCollection().Find(ctx, filter)
	if err != nil {
		errCh <- err
		return
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var doc models.Document
		if err := cur.Decode(&doc); err != nil {
			errCh <- err
			continue
		}
		docsCh <- doc
	}
}
