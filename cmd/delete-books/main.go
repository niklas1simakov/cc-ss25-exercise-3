package main

import (
	"context"
	"log"
	"net/http"

	"github.com/CAPS-Cloud/exercises/common/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	port = "3004"
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

	e.DELETE("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id")

		// Check if the book exists, otherwise not possible to delete
		existingBook, err := coll.FindOne(context.TODO(), bson.M{"id": id}).Raw()
		if err != nil && err != mongo.ErrNoDocuments {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		if existingBook == nil {
			return c.JSON(http.StatusNotFound, "Book not found")
		}

		// Delete the book
		result, err := database.DeleteOneBook(coll, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, result)
	})

	log.Printf("Delete-Books service started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
