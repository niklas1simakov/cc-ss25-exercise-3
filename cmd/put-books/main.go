package main

import (
	"log"
	"net/http"

	"github.com/CAPS-Cloud/exercises/common/database"
	"github.com/CAPS-Cloud/exercises/common/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	dbName     = "bookstore"
	collecName = "books"
	port       = "3003"
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

	e.PUT("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id")
		var book models.BookStore
		if err := c.Bind(&book); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid book format"})
		}

		result, err := database.UpdateOneBook(coll, id, book)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update book"})
		}

		return c.JSON(http.StatusOK, result)
	})

	log.Printf("Put-Books service started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
