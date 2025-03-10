package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"

	"github.com/fiantyogalihp/agn-cetak-toolbox/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

//go:embed templates/*
var templatesFS embed.FS

//go:embed screens/*.json
var jsonFiles embed.FS

func main() {
	// Set up Fiber with the HTML template engine
	engine := html.NewFileSystem(http.FS(templatesFS), ".html")

	engine.AddFunc("toJson", func(v interface{}) (string, error) {
		bytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// MIDDLEWARE
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} |${ip}| ${status} - ${method} ${path}\n",
		TimeFormat: "02 Jan 2006 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	// Serve static files (e.g., CSS/JS)
	app.Static("/static", "./static")

	// SET ROUTES
	routers.SetRouters(app, jsonFiles)

	// Start the server
	log.Fatal(app.Listen(":36530"))
}
