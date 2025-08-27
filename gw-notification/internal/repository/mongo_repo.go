package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoRepository(uri, dbName, collection string) (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	col := client.Database(dbName).Collection(collection)

	_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]int{"event_id": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return &MongoRepository{
		client:     client,
		collection: client.Database(dbName).Collection(collection),
	}, nil
}

func (r *MongoRepository) SaveEvents(ctx context.Context, event interface{}) error {
	_, err := r.collection.InsertOne(ctx, event)
	return err
}

func (s *MongoRepository) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}
