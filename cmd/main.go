package main

import (
	"artela-service/internal/config"
	"artela-service/internal/entity"
	"artela-service/internal/handler"
	"artela-service/internal/repository"
	"artela-service/internal/service"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found, using system env")
	}

	// 2. Init Database
	db := config.NewDatabase()
	db.AutoMigrate(&entity.Invitation{}, &entity.GalleryImage{})

	// 3. Dependency Injection
	invRepo := repository.NewInvitationRepository(db)
	invService := service.NewInvitationService(invRepo)
	invHandler := handler.NewInvitationHandler(invService)
	healthHandler := handler.NewHealthHandler(db)

	// 4. Setup Fiber
	app := fiber.New()
	app.Use(cors.New())

	// 5. Routes
	app.Get("/health", healthHandler.Check)

	api := app.Group("/api")
	api.Get("/invitation/:slug", invHandler.GetInvitation)
	api.Post("/admin/create", invHandler.CreateInvitation)

	// 6. Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}