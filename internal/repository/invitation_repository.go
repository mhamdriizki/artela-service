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
	
	CreateGallery(images []entity.GalleryImage) error
	FindGalleryImageByID(id string) (*entity.GalleryImage, error) // ID string (UUID)
	DeleteGalleryImage(id string) error
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
	err := r.db.Select("slug, couple_name, theme, created_at").Find(&invs).Error
	return invs, err
}

func (r *invitationRepository) FindBySlug(slug string) (*entity.Invitation, error) {
	var invitation entity.Invitation
	err := r.db.Preload("Gallery").Where("slug = ?", slug).First(&invitation).Error
	return &invitation, err
}

func (r *invitationRepository) Update(invitation *entity.Invitation) error {
	return r.db.Save(invitation).Error
}

func (r *invitationRepository) Delete(slug string) error {
	// Unscoped() = Hard Delete (Hapus permanen dari DB)
	return r.db.Unscoped().Where("slug = ?", slug).Delete(&entity.Invitation{}).Error
}

func (r *invitationRepository) CreateGallery(images []entity.GalleryImage) error {
	return r.db.Create(&images).Error
}

func (r *invitationRepository) FindGalleryImageByID(id string) (*entity.GalleryImage, error) {
	var img entity.GalleryImage
	err := r.db.First(&img, "id = ?", id).Error
	return &img, err
}

func (r *invitationRepository) DeleteGalleryImage(id string) error {
	// Unscoped() = Hard Delete
	return r.db.Unscoped().Delete(&entity.GalleryImage{}, "id = ?", id).Error
}