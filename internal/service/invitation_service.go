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
}

// Implementation
type invitationService struct {
	repo repository.InvitationRepository
}

// Dependency Injection via Constructor
func NewInvitationService(repo repository.InvitationRepository) InvitationService {
	return &invitationService{repo: repo}
}

func (s *invitationService) GetInvitation(slug string) (*entity.InvitationResponse, error) {
	// 1. Ambil data dari repository
	inv, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("undangan tidak ditemukan")
	}

	// 2. Transformasi ke format JSON Angular (Business Logic)
	galleryUrls := []string{}
	for _, img := range inv.GalleryImages {
		galleryUrls = append(galleryUrls, img.Url)
	}

	response := &entity.InvitationResponse{
		Slug:          inv.Slug,
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