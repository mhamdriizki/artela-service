package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseEntity menggantikan gorm.Model untuk support UUID
type BaseEntity struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key;" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Kita hilangkan DeletedAt di sini agar default-nya Hard Delete.
	// Jika butuh Soft Delete, tambahkan manual di struct spesifik.
}

// Hook otomatis generate UUID sebelum create
func (base *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return
}