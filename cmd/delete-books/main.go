package main

import (
	"log"
	"net/http"

	"github.com/CAPS-Cloud/exercises/common/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	dbName     = "bookstore"
	collecName = "books"
	port       = "3004"
)

func main() {
	client, err := database.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	coll, err := database.PrepareDatabase(client, dbName, collecName)
	if err != nil {
		log.Fatalf("failed to prepare database: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.DELETE("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id")
		result, err := database.DeleteOneBook(coll, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete book"})
		}
		if result.DeletedCount == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "book not found"})
		}

		return c.JSON(http.StatusOK, result)
	})

	log.Printf("Delete-Books service started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
