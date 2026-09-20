package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAdapter "mongodb-api/internal/adapters/http"
	"mongodb-api/internal/adapters/repository"
	"mongodb-api/internal/core/services"
	"mongodb-api/internal/infrastructure/config"
	"mongodb-api/internal/infrastructure/db"
	"mongodb-api/internal/infrastructure/logger"
)

func main() {
	// 1. Infrastructure: Load environment configuration
	cfg := config.Load()

	// 2. Infrastructure: Connect to MongoDB
	mongoClient, err := db.NewMongoClient(cfg.MongoURI)
	if err != nil {
		logger.New(logger.Warning).Print("Iniciando servidor sin conexión activa a MongoDB (comprueba MONGO_URI en .env)")
	}

	// 3. Adapters (Outbound): Initialize MongoDB Repository
	documentRepo := repository.NewMongoRepository(mongoClient)

	// 4. Core Domain: Initialize Use Cases / Business Logic
	documentService := services.NewDocumentService(documentRepo)

	// 5. Adapters (Inbound): Initialize HTTP Handlers & Router
	documentHandler := httpAdapter.NewDocumentHandler(documentService)
	app := httpAdapter.NewRouter(documentHandler)

	addr := fmt.Sprintf(":%s", cfg.Port)

	// 6. Start server
	go func() {
		logger.New(logger.Success).Print(fmt.Sprintf("🚀 Servidor iniciado en http://localhost:%s (Arquitectura Hexagonal + Fiber)", cfg.Port))
		if err := app.Listen(addr); err != nil {
			logger.New(logger.Error).Print(fmt.Sprintf("Error en servidor Fiber: %v", err))
		}
	}()

	// 7. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.New(logger.System).Print("Apagando servidor...")

	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		logger.New(logger.Error).Print(fmt.Sprintf("Apagado forzado del servidor: %v", err))
	}

	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mongoClient.Close(ctx)
	}

	logger.New(logger.Success).Print("Servidor finalizado limpiamente.")
}
