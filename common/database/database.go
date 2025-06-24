package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CAPS-Cloud/exercises/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("DATABASE_URI")
	if len(uri) == 0 {
		return nil, fmt.Errorf("failure to load env variable DATABASE_URI")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to create client for MongoDB: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB, please make sure the database is running: %w", err)
	}

	return client, nil
}

func PrepareDatabase(client *mongo.Client, dbName string, collecName string) (*mongo.Collection, error) {
	db := client.Database(dbName)

	// A more robust way to do this would be to check for a specific error,
	// but for this exercise, we'll just drop if it exists to ensure a clean state.
	db.Collection(collecName).Drop(context.TODO())

	cmd := bson.D{{Key: "create", Value: collecName}}
	var result bson.M
	if err := db.RunCommand(context.TODO(), cmd).Decode(&result); err != nil {
		log.Printf("Failed to create collection: %v", err)
		return nil, err
	}

	return db.Collection(collecName), nil
}

func PrepareData(coll *mongo.Collection) {
	startData := []models.BookStore{
		{
			ID:         "example1",
			BookName:   "The Vortex",
			BookAuthor: "José Eustasio Rivera",
			BookPages:  "292",
			BookYear:   "1924",
		},
		{
			ID:         "example2",
			BookName:   "Frankenstein",
			BookAuthor: "Mary Shelley",
			BookPages:  "280",
			BookYear:   "1818",
		},
		{
			ID:         "example3",
			BookName:   "The Black Cat",
			BookAuthor: "Edgar Allan Poe",
			BookPages:  "280",
			BookYear:   "1843",
		},
	}

	for _, book := range startData {
		_, err := coll.InsertOne(context.TODO(), book)
		if err != nil {
			// In a real app, you'd want more sophisticated error handling.
			log.Printf("Could not insert book %s: %v", book.BookName, err)
		}
	}
}

func FindAllBooks(coll *mongo.Collection) ([]models.BookStore, error) {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	var results []models.BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func FindAllAuthors(coll *mongo.Collection) ([]map[string]interface{}, error) {
	// Using an aggregation pipeline is more efficient for getting distinct values.
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$author"}}}},
		{{Key: "$project", Value: bson.D{{Key: "BookAuthor", Value: "$_id"}, {Key: "_id", Value: 0}}}},
	}
	cursor, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func FindAllYears(coll *mongo.Collection) ([]map[string]interface{}, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$year"}}}},
		{{Key: "$project", Value: bson.D{{Key: "BookYear", Value: "$_id"}, {Key: "_id", Value: 0}}}},
	}
	cursor, err := coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func InsertOneBook(coll *mongo.Collection, book models.BookStore) (*mongo.InsertOneResult, error) {
	return coll.InsertOne(context.TODO(), book)
}

func UpdateOneBook(coll *mongo.Collection, id string, book models.BookStore) (*mongo.UpdateResult, error) {
	filter := bson.D{{Key: "id", Value: id}}
	update := bson.D{{Key: "$set", Value: book}}
	return coll.UpdateOne(context.TODO(), filter, update)
}

func DeleteOneBook(coll *mongo.Collection, id string) (*mongo.DeleteResult, error) {
	filter := bson.D{{Key: "id", Value: id}}
	return coll.DeleteOne(context.TODO(), filter)
}
