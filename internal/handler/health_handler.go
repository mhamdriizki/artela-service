package handler

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	// 1. Cek Koneksi Database (Ping)
	sqlDB, err := h.db.DB()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "error", 
			"message": "Database instance error",
		})
	}

	if err := sqlDB.Ping(); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status": "error", 
			"message": "Database connection lost",
		})
	}

	// 2. Return OK jika semua aman
	return c.Status(200).JSON(fiber.Map{
		"status": "UP",
		"database": "connected",
		"service": "artela-api",
	})
}