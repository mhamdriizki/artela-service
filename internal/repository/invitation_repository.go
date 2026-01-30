package repository

import (
	"artela-service/internal/entity"

	"gorm.io/gorm"
)

type InvitationRepository interface {
	FindAll() ([]entity.Invitation, error)
	FindBySlug(slug string) (*entity.Invitation, error)
	FindGalleryImageByID(id string) (*entity.GalleryImage, error)
	
	Create(invitation *entity.Invitation) error
	Update(invitation *entity.Invitation) error
	UpdateCouplePhotos(slug string, groomPhoto string, bridePhoto string) error
	Delete(slug string) error
	
	CreateGallery(images []entity.GalleryImage) error
	DeleteGalleryImage(id string) error
	
	// METHOD BARU
	CreateGuestbook(guestbook *entity.Guestbook) error
	CreateRSVP(rsvp *entity.RSVP) error
}

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) FindAll() ([]entity.Invitation, error) {
	var invitations []entity.Invitation
	err := r.db.Order("created_at desc").Find(&invitations).Error
	return invitations, err
}

func (r *invitationRepository) FindBySlug(slug string) (*entity.Invitation, error) {
	var invitation entity.Invitation
	
	// UPDATE: Tambahkan Preload("Guestbooks") diurutkan dari yang terbaru
	err := r.db.Preload("Gallery").
		Preload("Guestbooks", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Preload("RSVPs").
		Where("slug = ?", slug).
		First(&invitation).Error
		
	return &invitation, err
}

func (r *invitationRepository) Create(invitation *entity.Invitation) error {
	return r.db.Create(invitation).Error
}

func (r *invitationRepository) Update(invitation *entity.Invitation) error {
	return r.db.Save(invitation).Error
}

func (r *invitationRepository) UpdateCouplePhotos(slug string, groomPhoto string, bridePhoto string) error {
	updates := map[string]interface{}{}
	if groomPhoto != "" {
		updates["groom_photo"] = groomPhoto
	}
	if bridePhoto != "" {
		updates["bride_photo"] = bridePhoto
	}
	return r.db.Model(&entity.Invitation{}).Where("slug = ?", slug).Updates(updates).Error
}

func (r *invitationRepository) Delete(slug string) error {
	return r.db.Where("slug = ?", slug).Delete(&entity.Invitation{}).Error
}

func (r *invitationRepository) CreateGallery(images []entity.GalleryImage) error {
	return r.db.Create(&images).Error
}

func (r *invitationRepository) FindGalleryImageByID(id string) (*entity.GalleryImage, error) {
	var img entity.GalleryImage
	err := r.db.Where("id = ?", id).First(&img).Error
	return &img, err
}

func (r *invitationRepository) DeleteGalleryImage(id string) error {
	return r.db.Where("id = ?", id).Delete(&entity.GalleryImage{}).Error
}

// IMPLEMENTASI BARU
func (r *invitationRepository) CreateGuestbook(guestbook *entity.Guestbook) error {
	return r.db.Create(guestbook).Error
}

func (r *invitationRepository) CreateRSVP(rsvp *entity.RSVP) error {
	return r.db.Create(rsvp).Error
}