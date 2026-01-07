package repository

import (
	"artela-service/internal/entity"

	"gorm.io/gorm"
)

// Interface (Kontrak)
type InvitationRepository interface {
	FindBySlug(slug string) (*entity.Invitation, error)
	Create(invitation *entity.Invitation) error
	AddGallery(inv *entity.Invitation, images []entity.GalleryImage) error
}

// Implementation
type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) FindBySlug(slug string) (*entity.Invitation, error) {
	var inv entity.Invitation
	err := r.db.Preload("GalleryImages").Where("slug = ?", slug).First(&inv).Error
	return &inv, err
}

func (r *invitationRepository) Create(invitation *entity.Invitation) error {
	return r.db.Create(invitation).Error
}

func (r *invitationRepository) AddGallery(inv *entity.Invitation, images []entity.GalleryImage) error {
    return r.db.Model(inv).Association("GalleryImages").Append(images)
}