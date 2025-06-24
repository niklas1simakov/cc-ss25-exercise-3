package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Defines a "model" that we can use to communicate with the
// frontend or the database
// More on these "tags" like `bson:"_id,omitempty"`: https://go.dev/wiki/Well-known-struct-tags
type BookStore struct {
	MongoID     primitive.ObjectID `bson:"_id,omitempty"`
	ID          string             `bson:"id" json:"id" form:"id"`
	BookName    string             `bson:"title" json:"title" form:"title"`
	BookAuthor  string             `bson:"author" json:"author" form:"author"`
	BookEdition string             `bson:"edition,omitempty" json:"edition" form:"edition"`
	BookPages   string             `bson:"pages,omitempty" json:"pages" form:"pages"`
	BookYear    string             `bson:"year,omitempty" json:"year" form:"year"`
}
