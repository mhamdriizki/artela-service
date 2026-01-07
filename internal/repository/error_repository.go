package repository

import (
	"artela-service/internal/entity"

	"gorm.io/gorm"
)

type ErrorRepository interface {
	FindByCode(code string) entity.ErrorReference
}

type errorRepository struct {
	db *gorm.DB
}

func NewErrorRepository(db *gorm.DB) ErrorRepository {
	return &errorRepository{db: db}
}

func (r *errorRepository) FindByCode(code string) entity.ErrorReference {
	var errRef entity.ErrorReference
	// Jika tidak ketemu, return default error
	if result := r.db.First(&errRef, "code = ?", code); result.Error != nil {
		return entity.ErrorReference{
			Code:      "ART-99-999",
			MessageEN: "System Error (Unknown Code)",
			MessageID: "Kesalahan Sistem (Kode Tidak Dikenal)",
		}
	}
	return errRef
}