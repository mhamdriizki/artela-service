package service

import (
	"artela-service/internal/entity"
	"artela-service/internal/repository"
	"errors"
)

// Interface
type InvitationService interface {
	GetInvitation(slug string) (*entity.InvitationResponse, error)
	CreateInvitation(req *entity.Invitation) error
	AddGalleryImages(slug string, imageUrls []string) error
	UpdateInvitation(slug string, req *entity.Invitation) error // Baru
	DeleteInvitation(slug string) error                         // Baru
}

// Implementation
type invitationService struct {
	repo repository.InvitationRepository
}

func NewInvitationService(repo repository.InvitationRepository) InvitationService {
	return &invitationService{repo: repo}
}

func (s *invitationService) GetInvitation(slug string) (*entity.InvitationResponse, error) {
	// 1. Ambil data
	inv, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("undangan tidak ditemukan")
	}

	// 2. Transformasi ke Response
	galleryUrls := []string{}
	for _, img := range inv.GalleryImages {
		galleryUrls = append(galleryUrls, img.Url)
	}

	response := &entity.InvitationResponse{
		Slug:          inv.Slug,
		Theme:         inv.Theme, // Mapping Theme
		CoupleName:    inv.CoupleName,
		GroomName:     inv.GroomName,
		GroomPhoto:    inv.GroomPhoto,
		BrideName:     inv.BrideName,
		BridePhoto:    inv.BridePhoto,
		YoutubeUrl:    inv.YoutubeUrl,
		Gallery:       galleryUrls,
		EventDetails: entity.EventDetails{
			Date:     inv.EventDate,
			Location: inv.EventLocation,
			Address:  inv.EventAddress,
			MapUrl:   inv.MapUrl,
		},
	}

	return response, nil
}

func (s *invitationService) CreateInvitation(req *entity.Invitation) error {
	return s.repo.Create(req)
}

func (s *invitationService) AddGalleryImages(slug string, imageUrls []string) error {
	inv, err := s.repo.FindBySlug(slug)
	if err != nil {
		return err
	}

	var gallery []entity.GalleryImage
	for _, url := range imageUrls {
		gallery = append(gallery, entity.GalleryImage{
			InvitationID: inv.ID,
			Url:          url,
		})
	}
	return s.repo.AddGallery(inv, gallery)
}

// Logic Update
func (s *invitationService) UpdateInvitation(slug string, req *entity.Invitation) error {
	// 1. Cek data lama
	existing, err := s.repo.FindBySlug(slug)
	if err != nil {
		return errors.New("data tidak ditemukan")
	}

	// 2. Assign ID lama ke struct baru (agar terhitung update)
	req.ID = existing.ID
	req.CreatedAt = existing.CreatedAt
	
	// Handle jika slug kosong, pakai yang lama
	if req.Slug == "" {
		req.Slug = existing.Slug
	}

	return s.repo.Update(req)
}

// Logic Delete
func (s *invitationService) DeleteInvitation(slug string) error {
	existing, err := s.repo.FindBySlug(slug)
	if err != nil {
		return errors.New("data tidak ditemukan")
	}
	return s.repo.Delete(existing)
}