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

// Fungsi Seeder Sederhana
func seedErrorCodes(db *gorm.DB) {
	codes := []entity.ErrorReference{
		// SUCCESS (200) - Prefix 00
		{Code: "ART-00-000", MessageEN: "Success", MessageID: "Berhasil"},
		{Code: "ART-00-001", MessageEN: "Created Successfully", MessageID: "Data Berhasil Dibuat"},
		
		// CLIENT ERRORS (4xx) - Prefix 98
		{Code: "ART-98-001", MessageEN: "Invalid Input Data", MessageID: "Data Input Tidak Valid"},
		{Code: "ART-98-002", MessageEN: "Missing Required Fields", MessageID: "Data Wajib Belum Diisi"},
		{Code: "ART-98-003", MessageEN: "Slug Already Exists", MessageID: "Link Undangan Sudah Digunakan"},
		{Code: "ART-98-004", MessageEN: "Data Not Found", MessageID: "Data Tidak Ditemukan"},
		{Code: "ART-98-005", MessageEN: "Max 5 files allowed", MessageID: "Maksimal upload 5 foto sekaligus"},
    {Code: "ART-98-006", MessageEN: "Invalid file type (JPG/PNG only)", MessageID: "Format file salah (Hanya JPG/PNG)"},
		{Code: "ART-98-100", MessageEN: "Unauthorized", MessageID: "Akses Ditolak"},
		
		// SERVER ERRORS (5xx) - Prefix 99
		{Code: "ART-99-000", MessageEN: "Request Timeout", MessageID: "Waktu Permintaan Habis"},
		{Code: "ART-99-001", MessageEN: "Database Error", MessageID: "Kesalahan Database"},
		{Code: "ART-99-999", MessageEN: "Internal Server Error", MessageID: "Terjadi Kesalahan Sistem"},
	}

	for _, c := range codes {
		// Gunakan Clauses(clause.OnConflict) jika ingin update data jika code sudah ada
		// Atau FirstOrCreate seperti sebelumnya
		db.FirstOrCreate(&c, entity.ErrorReference{Code: c.Code})
	}
	log.Println("✅ Error Codes Dictionary seeded completely!")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: using system env")
	}

	// 1. Init Database
	db := config.NewDatabase()
	// Migrate tabel baru ErrorReference
	db.AutoMigrate(&entity.Invitation{}, &entity.GalleryImage{}, &entity.Guestbook{}, &entity.RSVP{}, &entity.ErrorReference{})

	// 2. Run Seeder
	seedErrorCodes(db)

	// 3. Dependency Injection
	// Repositories
	invRepo := repository.NewInvitationRepository(db)
	errRepo := repository.NewErrorRepository(db) // Repo baru

	// Services
	invService := service.NewInvitationService(invRepo)

	// Handlers
	// Masukkan errRepo ke InvitationHandler
	invHandler := handler.NewInvitationHandler(invService, errRepo) 
	healthHandler := handler.NewHealthHandler(db)

	// 4. Fiber App
	app := fiber.New()
	app.Use(cors.New())

	app.Static("/uploads", "./public/uploads")

	// Health check (RAW response, tidak pakai schema wrapper)
	app.Get("/health", healthHandler.Check)

	api := app.Group("/api")
	api.Get("/invitation/:slug", invHandler.GetInvitation)
	api.Post("/admin/create", invHandler.CreateInvitation)
	api.Post("/invitation/:slug/gallery", invHandler.UploadGallery)

	// Start
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}