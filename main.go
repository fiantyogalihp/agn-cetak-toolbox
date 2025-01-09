package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/fiantyogalihp/dynamic-json-parsing-struct/routers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

//go:embed templates/*
var templatesFS embed.FS

func main() {
	// Set up Fiber with the HTML template engine
	engine := html.NewFileSystem(http.FS(templatesFS), ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Serve static files (e.g., CSS/JS)
	app.Static("/static", "./static")

	routers.SetRouters(app)

	// Start the server
	log.Fatal(app.Listen(":36530"))
}
