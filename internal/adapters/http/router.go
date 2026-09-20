package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewRouter(handler *DocumentHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               "MongoDB REST API (Hexagonal - Fiber)",
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))

	// Health & Monitor
	app.Get("/health", handler.HealthCheck)
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MongoDB REST API - Server Metrics"}))
	app.Get("/monitor", monitor.New(monitor.Config{Title: "MongoDB REST API - Server Monitor"}))

	// API CRUD routes
	api := app.Group("/api")
	{
		api.Get("/:database/:collection", handler.GetAll)
		api.Get("/:database/:collection/:id", handler.GetByID)
		api.Post("/:database/:collection", handler.Create)
		api.Put("/:database/:collection/:id", handler.Replace)
		api.Patch("/:database/:collection/:id", handler.UpdatePartial)
		api.Delete("/:database/:collection/:id", handler.Delete)
		api.Delete("/:database/:collection", handler.Delete)
	}

	return app
}
