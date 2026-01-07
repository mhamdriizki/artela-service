package utils

import (
	"artela-service/internal/entity"
	"artela-service/internal/repository"

	"github.com/gofiber/fiber/v2"
)

// Helper untuk format sukses/gagal standar
func BuildResponse(ctx *fiber.Ctx, errorRepo repository.ErrorRepository, code string, data interface{}) error {
	errRef := errorRepo.FindByCode(code)

	httpStatus := 200 // Default Success

	// Logic Mapping Prefix ke HTTP Status
	prefix := code[0:7] // Ambil 7 karakter pertama "ART-xx-"

	if prefix == "ART-00-" {
		httpStatus = 200
	} else if prefix == "ART-98-" {
		// Default Client Error
		httpStatus = 400 
		
		// Case khusus jika Not Found
		if code == "ART-98-004" {
			httpStatus = 404
		}
		// Case khusus Unauthorized
		if code == "ART-98-100" || code == "ART-98-101" {
			httpStatus = 401
		}
	} else if prefix == "ART-99-" {
		// Server Error
		httpStatus = 500
		if code == "ART-99-000" {
			httpStatus = 408 // Timeout
		}
	}

	response := entity.APIResponse{
		ErrorSchema: entity.ErrorSchema{
			ErrorCode: errRef.Code,
			ErrorMessage: entity.ErrorMessage{
				English:    errRef.MessageEN,
				Indonesian: errRef.MessageID,
			},
		},
		OutputSchema: data,
	}

	return ctx.Status(httpStatus).JSON(response)
}