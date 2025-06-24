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
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid book format"})
		}

		result, err := database.InsertOneBook(coll, book)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create book"})
		}

		return c.JSON(http.StatusCreated, result)
	})

	log.Printf("Post-Books service started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
