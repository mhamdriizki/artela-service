package main

import (
	"artela-service/internal/config"
	"artela-service/internal/entity"
	"artela-service/internal/handler"
	"artela-service/internal/middleware"
	"artela-service/internal/repository"
	"artela-service/internal/service"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ... (Func seedErrorCodes & seedAdmin tetap sama, tidak saya tulis ulang untuk hemat space) ...
func seedErrorCodes(db *gorm.DB) {
	// ... (Code existing) ...
    // Pastikan code ART-98-001 dll ada
    db.FirstOrCreate(&entity.ErrorReference{Code: "ART-00-000", MessageEN: "Success", MessageID: "Berhasil"})
    // ... tambahkan yang lain jika perlu ...
}

func seedAdmin(db *gorm.DB) {
	// ... (Code existing) ...
    var count int64
	db.Model(&entity.User{}).Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("Langit1105"), bcrypt.DefaultCost)
		db.Create(&entity.User{Username: "admin01", Password: string(hash), Role: "admin"})
	}
}

func main() {
	godotenv.Load()

	// Ensure Upload Folder
	if _, err := os.Stat("./public/uploads"); os.IsNotExist(err) {
		os.MkdirAll("./public/uploads", 0755)
	}

	db := config.NewDatabase()
	
	// MIGRATE: Tambahkan Guestbook
	db.AutoMigrate(
		&entity.Invitation{},
		&entity.GalleryImage{},
		&entity.Guestbook{}, // <--- NEW
		&entity.RSVP{},
		&entity.ErrorReference{},
		&entity.User{},
	)

	seedErrorCodes(db)
	seedAdmin(db)

	// Dependency Injection
	invRepo := repository.NewInvitationRepository(db)
	errRepo := repository.NewErrorRepository(db)
	invService := service.NewInvitationService(invRepo)
	invHandler := handler.NewInvitationHandler(invService, errRepo)
	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(db)

	app := fiber.New(fiber.Config{ BodyLimit: 30 * 1024 * 1024 })
	app.Use(cors.New())
	app.Static("/uploads", "./public/uploads")

	app.Get("/health", healthHandler.Check)
	
	api := app.Group("/api")
	api.Post("/login", authHandler.Login)
	api.Post("/logout", authHandler.Logout)
	
	// PUBLIC ROUTES
	api.Get("/invitation/:slug", invHandler.GetInvitation)
	api.Post("/invitation/:slug/gallery", invHandler.UploadGallery) // Gallery Existing (Public?)
	api.Post("/invitation/:slug/guestbook", invHandler.CreateGuestbook) // <--- ROUTE BARU
	api.Post("/invitation/:slug/rsvp", invHandler.CreateRSVP)

	// ADMIN ROUTES
	admin := api.Group("/admin")
	admin.Use(middleware.Protected())
	
	admin.Get("/invitations", invHandler.GetAllInvitations)
	admin.Post("/create", invHandler.CreateInvitation)
	admin.Put("/invitation/:slug", invHandler.UpdateInvitation)
	admin.Post("/invitation/:slug/upload-couple", invHandler.UploadCouplePhotos)
	admin.Delete("/invitation/:slug", invHandler.DeleteInvitation)
	admin.Delete("/gallery/:id", invHandler.DeleteGalleryImage)

	port := os.Getenv("APP_PORT")
	if port == "" { port = "3000" }
	log.Fatal(app.Listen(":" + port))
}