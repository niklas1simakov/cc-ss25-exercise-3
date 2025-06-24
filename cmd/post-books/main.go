package main

import (
	"context"
	"log"
	"net/http"

	"github.com/CAPS-Cloud/exercises/common/database"
	"github.com/CAPS-Cloud/exercises/common/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	port = "3002"
)

func main() {
	client := database.Connect()

	coll, err := database.PrepareDatabase(client)
	if err != nil {
		log.Fatalf("failed to prepare database: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/books", func(c echo.Context) error {
		var book models.BookStore
		if err := c.Bind(&book); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if book.ID == "" || book.BookName == "" || book.BookAuthor == "" {
			return c.JSON(http.StatusBadRequest, "Missing required fields")
		}

		// Check if the book already exists
		newBookId := bson.M{"id": book.ID}
		existingBook, err := coll.FindOne(context.TODO(), newBookId).Raw()
		if err != nil && err != mongo.ErrNoDocuments {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if existingBook != nil {
			return c.JSON(http.StatusConflict, "Book already exists")
		}

		// Add the book to the database
		newBook, err := database.InsertOneBook(coll, book)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, newBook)
	})

	log.Printf("Post-Books service started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
