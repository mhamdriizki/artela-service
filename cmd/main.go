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
	"gorm.io/gorm"
)

func seedErrorCodes(db *gorm.DB) {
	codes := []entity.ErrorReference{
		{Code: "ART-00-000", MessageEN: "Success", MessageID: "Berhasil"},
		{Code: "ART-00-001", MessageEN: "Created Successfully", MessageID: "Data Berhasil Dibuat"},
		{Code: "ART-00-002", MessageEN: "Updated Successfully", MessageID: "Data Berhasil Diperbarui"}, // Baru
		{Code: "ART-00-003", MessageEN: "Deleted Successfully", MessageID: "Data Berhasil Dihapus"},    // Baru
		{Code: "ART-98-001", MessageEN: "Invalid Input Data", MessageID: "Data Input Tidak Valid"},
		{Code: "ART-98-004", MessageEN: "Data Not Found", MessageID: "Data Tidak Ditemukan"},
		{Code: "ART-98-005", MessageEN: "Max 5 files allowed", MessageID: "Maksimal upload 5 foto sekaligus"},
		{Code: "ART-98-006", MessageEN: "Invalid file type", MessageID: "Format file salah (Hanya JPG/PNG)"},
		{Code: "ART-99-002", MessageEN: "Database Error", MessageID: "Kesalahan Database"},
		{Code: "ART-99-999", MessageEN: "Internal Server Error", MessageID: "Terjadi Kesalahan Sistem"},
	}

	for _, c := range codes {
		db.FirstOrCreate(&c, entity.ErrorReference{Code: c.Code})
	}
	log.Println("✅ Error Codes seeded")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: using system env")
	}

	// 1. Init DB
	db := config.NewDatabase()
	db.AutoMigrate(&entity.Invitation{}, &entity.GalleryImage{}, &entity.Guestbook{}, &entity.RSVP{}, &entity.ErrorReference{})

	// 2. Seed
	seedErrorCodes(db)

	// 3. DI
	invRepo := repository.NewInvitationRepository(db)
	errRepo := repository.NewErrorRepository(db)
	invService := service.NewInvitationService(invRepo)
	invHandler := handler.NewInvitationHandler(invService, errRepo)
	healthHandler := handler.NewHealthHandler(db)

	// 4. Fiber
	app := fiber.New()
	app.Use(cors.New())
	app.Static("/uploads", "./public/uploads")

	app.Get("/health", healthHandler.Check)

	api := app.Group("/api")
	api.Get("/invitation/:slug", invHandler.GetInvitation)
	api.Post("/invitation/:slug/gallery", invHandler.UploadGallery)

	admin := api.Group("/admin")
	admin.Post("/create", invHandler.CreateInvitation)
	admin.Put("/invitation/:slug", invHandler.UpdateInvitation)    // Endpoint Baru
	admin.Delete("/invitation/:slug", invHandler.DeleteInvitation) // Endpoint Baru

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}