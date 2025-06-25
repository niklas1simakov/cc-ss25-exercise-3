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

func Connect() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("DATABASE_URI")
	if len(uri) == 0 {
		fmt.Printf("failure to load env variable\n")
		os.Exit(1)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Printf("failed to create client for MongoDB\n")
		os.Exit(1)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("failed to connect to MongoDB, please make sure the database is running\n")
		os.Exit(1)
	}

	return client
}

func PrepareDatabase(client *mongo.Client) (*mongo.Collection, error) {
	dbName := "exercise-3"
	collecName := "information"

	db := client.Database(dbName)

	// A more robust way to do this would be to check for a specific error,
	// but for this exercise, we'll just drop if it exists to ensure a clean state.
	// db.Collection(collecName).Drop(context.TODO())

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
			ID:          "example1",
			BookName:    "The Vortex",
			BookAuthor:  "José Eustasio Rivera",
			BookEdition: "958-30-0804-4",
			BookPages:   "292",
			BookYear:    "1924",
		},
		{
			ID:          "example2",
			BookName:    "Frankenstein",
			BookAuthor:  "Mary Shelley",
			BookEdition: "978-3-649-64609-9",
			BookPages:   "280",
			BookYear:    "1818",
		},
		{
			ID:          "example3",
			BookName:    "The Black Cat",
			BookAuthor:  "Edgar Allan Poe",
			BookEdition: "978-3-99168-238-7",
			BookPages:   "280",
			BookYear:    "1843",
		},
	}

	// This syntax helps us iterate over arrays. It behaves similar to Python
	// However, range always returns a tuple: (idx, elem). You can ignore the idx
	// by using _.
	// In the topic of function returns: sadly, there is no standard on return types from function. Most functions
	// return a tuple with (res, err), but this is not granted. Some functions
	// might return a ret value that includes res and the err, others might have
	// an out parameter.
	for _, book := range startData {
		filter := bson.M{"id": book.ID}
		cursor, err := coll.Find(context.TODO(), filter)
		var results []models.BookStore
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		if len(results) > 0 {
			log.Printf("Book with ID %s already exists, skipping insertion.", book.ID)
		} else {
			result, err := coll.InsertOne(context.TODO(), book)
			if err != nil {
				panic(err)
			} else {
				fmt.Printf("%+v\n", result)
			}
		}
	}
}

// Generic method to perform "SELECT * FROM BOOKS" (if this was SQL, which
// it is not :D ), and then we convert it into an array of map. In Golang, you
// define a map by writing map[<key type>]<value type>{<key>:<value>}.
// interface{} is a special type in Golang, basically a wildcard...
func FindAllBooks(coll *mongo.Collection) []models.BookStore {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}
	var results []models.BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results
}

func FindAllAuthors(coll *mongo.Collection) []map[string]interface{} {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	var results []models.BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// Use a map to store distinct authors
	distinctAuthors := make(map[string]struct{})
	for _, res := range results {
		distinctAuthors[res.BookAuthor] = struct{}{}
	}

	var ret []map[string]interface{}
	// Convert distinct authors from map keys to a slice of maps
	for author := range distinctAuthors {
		ret = append(ret, map[string]interface{}{
			"BookAuthor": author,
		})
	}

	return ret
}

func FindAllYears(coll *mongo.Collection) []map[string]interface{} {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	var results []models.BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// Use a map to store distinct years
	distinctYears := make(map[string]struct{})
	for _, res := range results {
		distinctYears[res.BookYear] = struct{}{}
	}

	var ret []map[string]interface{}
	// Convert distinct years from map keys to a slice of maps
	for year := range distinctYears {
		ret = append(ret, map[string]interface{}{
			"BookYear": year,
		})
	}

	return ret
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
