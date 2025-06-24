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
	port = "3003"
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

	e.PUT("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id")
		var book models.BookStore
		if err := c.Bind(&book); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if book.BookName == "" || book.BookAuthor == "" {
			return c.JSON(http.StatusBadRequest, "Missing required fields")
		}

		// Check if the book exists, otherwise not possible to update
		existingBook, err := coll.FindOne(context.TODO(), bson.M{"id": id}).Raw()
		if err != nil && err != mongo.ErrNoDocuments {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if existingBook == nil {
			return c.JSON(http.StatusNotFound, "Book not found")
		}

		if book.ID == "" {
			book.ID = id
		}

		// Check if the book id is the same as the one in the body
		if book.ID != id {
			return c.JSON(http.StatusBadRequest, "Book ID in URL and body must match")
		}

		// Update the book
		result, err := database.UpdateOneBook(coll, id, book)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, result)
	})

	log.Printf("Put-Books service started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
