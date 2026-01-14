package repository

import (
	"artela-service/internal/entity"

	"gorm.io/gorm"
)

type InvitationRepository interface {
	Create(invitation *entity.Invitation) error
	FindAll() ([]entity.Invitation, error)
	FindBySlug(slug string) (*entity.Invitation, error)
	Update(invitation *entity.Invitation) error
	Delete(slug string) error
	
	// Gallery
	CreateGallery(images []entity.GalleryImage) error
	FindGalleryImageByID(id uint) (*entity.GalleryImage, error) // <-- BARU
	DeleteGalleryImage(id uint) error                          // <-- BARU
}

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) Create(invitation *entity.Invitation) error {
	return r.db.Create(invitation).Error
}

func (r *invitationRepository) FindAll() ([]entity.Invitation, error) {
	var invs []entity.Invitation
	// Select field penting saja agar query ringan
	err := r.db.Select("slug, couple_name, theme, created_at").Find(&invs).Error
	return invs, err
}

func (r *invitationRepository) FindBySlug(slug string) (*entity.Invitation, error) {
	var invitation entity.Invitation
	// Preload Gallery agar muncul saat detail dibuka
	err := r.db.Preload("Gallery").Where("slug = ?", slug).First(&invitation).Error
	return &invitation, err
}

func (r *invitationRepository) Update(invitation *entity.Invitation) error {
	return r.db.Save(invitation).Error
}

func (r *invitationRepository) Delete(slug string) error {
	return r.db.Where("slug = ?", slug).Delete(&entity.Invitation{}).Error
}

func (r *invitationRepository) CreateGallery(images []entity.GalleryImage) error {
	return r.db.Create(&images).Error
}

// --- IMPLEMENTASI BARU ---

func (r *invitationRepository) FindGalleryImageByID(id uint) (*entity.GalleryImage, error) {
	var img entity.GalleryImage
	err := r.db.First(&img, id).Error
	return &img, err
}

func (r *invitationRepository) DeleteGalleryImage(id uint) error {
	return r.db.Delete(&entity.GalleryImage{}, id).Error
}