package main

import (
	"log"
	"net/http"

	"github.com/CAPS-Cloud/exercises/common/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	port = "3001"
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

	e.GET("/api/books", func(c echo.Context) error {
		books := database.FindAllBooks(coll)
		return c.JSON(http.StatusOK, books)
	})

	log.Printf("Get-Books service started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
