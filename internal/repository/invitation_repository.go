package repository

import (
	"artela-service/internal/entity"

	"gorm.io/gorm"
)

// Interface
type InvitationRepository interface {
	FindBySlug(slug string) (*entity.Invitation, error)
	Create(invitation *entity.Invitation) error
	AddGallery(inv *entity.Invitation, images []entity.GalleryImage) error
	Update(invitation *entity.Invitation) error
	Delete(invitation *entity.Invitation) error
	FindAll() ([]entity.Invitation, error)
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

// Update Data
func (r *invitationRepository) Update(invitation *entity.Invitation) error {
	return r.db.Save(invitation).Error
}

// Delete Data
func (r *invitationRepository) Delete(invitation *entity.Invitation) error {
	return r.db.Delete(invitation).Error
}

func (r *invitationRepository) FindAll() ([]entity.Invitation, error) {
    var invs []entity.Invitation
    // Select field penting saja biar ringan
    err := r.db.Select("slug, couple_name, theme, created_at").Find(&invs).Error
    return invs, err
}