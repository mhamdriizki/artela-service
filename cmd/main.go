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

// --- SEEDERS ---

// 1. Seed Error Codes (Kamus Error)
func seedErrorCodes(db *gorm.DB) {
	codes := []entity.ErrorReference{
		{Code: "ART-00-000", MessageEN: "Success", MessageID: "Berhasil"},
		{Code: "ART-00-001", MessageEN: "Created Successfully", MessageID: "Data Berhasil Dibuat"},
		{Code: "ART-00-002", MessageEN: "Updated Successfully", MessageID: "Data Berhasil Diperbarui"},
		{Code: "ART-00-003", MessageEN: "Deleted Successfully", MessageID: "Data Berhasil Dihapus"},
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

// 2. Seed Admin User (Default Admin)
func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count == 0 {
		// Default: username 'admin''
		hash, _ := bcrypt.GenerateFromPassword([]byte("Langit1105"), bcrypt.DefaultCost)
		admin := entity.User{
			Username: "artela",
			Password: string(hash),
			Role:     "admin",
		}
		db.Create(&admin)
		log.Println("✅ Admin user seeded: admin / admin123")
	}
}

// --- MAIN PROGRAM ---

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("Info: using system env (no .env file found)")
	}

	// 2. Database Connection & Migration
	db := config.NewDatabase()
	
	// Auto Migrate semua Entity (termasuk User baru)
	err := db.AutoMigrate(
		&entity.Invitation{},
		&entity.GalleryImage{},
		&entity.Guestbook{},
		&entity.RSVP{},
		&entity.ErrorReference{},
		&entity.User{}, // Table User untuk Auth
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// 3. Run Seeders
	seedErrorCodes(db)
	seedAdmin(db)

	// 4. Dependency Injection (Wiring)
	
	// -- Repositories --
	invRepo := repository.NewInvitationRepository(db)
	errRepo := repository.NewErrorRepository(db)
	
	// -- Services --
	invService := service.NewInvitationService(invRepo)

	// -- Handlers --
	invHandler := handler.NewInvitationHandler(invService, errRepo)
	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(db) // Handler baru untuk Login

	// 5. Setup Fiber App
	app := fiber.New()
	
	// Middleware CORS (Agar bisa diakses Frontend berbeda domain/port)
	app.Use(cors.New())
	
	// Serve Static Files (Untuk Gambar Gallery)
	app.Static("/uploads", "./public/uploads")

	// 6. Define Routes
	
	// -- Public Routes --
	app.Get("/health", healthHandler.Check)
	
	api := app.Group("/api")
	api.Post("/login", authHandler.Login) // Login Endpoint
	
	// Public Invitation Routes
	api.Get("/invitation/:slug", invHandler.GetInvitation)
	api.Post("/invitation/:slug/gallery", invHandler.UploadGallery) // Upload Foto (Bisa dipindah ke admin jika mau strict)

	// -- Protected Admin Routes (Butuh Token Bearer) --
	admin := api.Group("/admin")
	admin.Use(middleware.Protected()) // Pasang Gembok Security 🔒

	// Dashboard: List Semua Undangan
	admin.Get("/invitations", invHandler.GetAllInvitations) 

	// CRUD Operations
	admin.Post("/create", invHandler.CreateInvitation)
	admin.Put("/invitation/:slug", invHandler.UpdateInvitation)
	admin.Delete("/invitation/:slug", invHandler.DeleteInvitation)

	// 7. Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("🚀 Server running on port " + port)
	log.Fatal(app.Listen(":" + port))
}