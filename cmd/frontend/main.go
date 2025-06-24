package main

import (
	"html/template"
	"io"
	"log"

	"github.com/CAPS-Cloud/exercises/common/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	port = "3000"
)

// Wraps the "Template" struct to associate a necessary method
// to determine the rendering procedure
type Template struct {
	tmpl *template.Template
}

// Preload the available templates for the view folder.
// This builds a local "database" of all available "blocks"
// to render upon request, i.e., replace the respective
// variable or expression.
// For more on templating, visit https://jinja.palletsprojects.com/en/3.0.x/templates/
// to get to know more about templating
// You can also read Golang's documentation on their templating
// https://pkg.go.dev/text/template
func loadTemplates() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

// Method definition of the required "Render" to be passed for the Rendering
// engine.
// Contraire to method declaration, such syntax defines methods for a given
// struct. "Interfaces" and "structs" can have methods associated with it.
// The difference lies that interfaces declare methods whether struct only
// implement them, i.e., only define them. Such differentiation is important
// for a compiler to ensure types provide implementations of such methods.
func (t *Template) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

func main() {
	client := database.Connect()

	coll, err := database.PrepareDatabase(client)
	if err != nil {
		log.Fatalf("failed to prepare database: %v", err)
	}

	// We only prepare the data in one service to avoid race conditions or duplication.
	// The frontend service is a good candidate for this.
	database.PrepareData(coll)

	e := echo.New()

	// Define our custom renderer
	e.Renderer = loadTemplates()

	// Log the requests. Please have a look at echo's documentation on more
	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/css", "css")

	// Endpoint definition. Here, we divided into two groups: top-level routes
	// starting with /, which usually serve webpages. For our RESTful endpoints,
	// we prefix the route with /api to indicate more information or resources
	// are available under such route.
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.GET("/books", func(c echo.Context) error {
		books := database.FindAllBooks(coll)
		return c.Render(200, "book-table", books)
	})

	e.GET("/authors", func(c echo.Context) error {
		authors := database.FindAllAuthors(coll)
		return c.Render(200, "authors-table", authors)
	})

	e.GET("/years", func(c echo.Context) error {
		years := database.FindAllYears(coll)
		return c.Render(200, "years-table", years)
	})

	e.GET("/search", func(c echo.Context) error {
		return c.Render(200, "search-bar", nil)
	})

	e.GET("/create", func(c echo.Context) error {
		return c.Render(200, "create-form", nil)
	})

	log.Printf("Frontend service started on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
