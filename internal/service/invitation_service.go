package service

import (
	"artela-service/internal/entity"
	"artela-service/internal/repository"
	"os"
	"path/filepath"
)

type InvitationService interface {
	GetAllInvitations() (*entity.InvitationListWrapper, error)
	GetInvitation(slug string) (*entity.Invitation, error)
	CreateInvitation(req *entity.Invitation) error
	UpdateInvitation(slug string, req *entity.Invitation) error
	DeleteInvitation(slug string) error
	UploadGallery(slug string, filenames []string) error
	DeleteGalleryImage(id string) error
}

type invitationService struct {
	repo repository.InvitationRepository
}

func NewInvitationService(repo repository.InvitationRepository) InvitationService {
	return &invitationService{repo: repo}
}

func (s *invitationService) GetAllInvitations() (*entity.InvitationListWrapper, error) {
	data, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var list []entity.InvitationListResponse
	for _, item := range data {
		list = append(list, entity.InvitationListResponse{
			Slug:        item.Slug,
			CoupleName:  item.CoupleName,
			Theme:       item.Theme,
			WeddingDate: item.WeddingDate.Format("2006-01-02"), // Format tanggal ke String
			CreatedAt:   item.CreatedAt.Format("2006-01-02"),
		})
	}
	
	if list == nil {
		list = []entity.InvitationListResponse{}
	}

	return &entity.InvitationListWrapper{Data: list}, nil
}

func (s *invitationService) GetInvitation(slug string) (*entity.Invitation, error) {
	return s.repo.FindBySlug(slug)
}

func (s *invitationService) CreateInvitation(req *entity.Invitation) error {
	return s.repo.Create(req)
}

func (s *invitationService) UpdateInvitation(slug string, req *entity.Invitation) error {
	existing, err := s.repo.FindBySlug(slug)
	if err != nil {
		return err
	}
	req.ID = existing.ID
	req.CreatedAt = existing.CreatedAt
	return s.repo.Update(req)
}

func (s *invitationService) DeleteInvitation(slug string) error {
	return s.repo.Delete(slug)
}

func (s *invitationService) UploadGallery(slug string, filenames []string) error {
	inv, err := s.repo.FindBySlug(slug)
	if err != nil {
		return err
	}

	var images []entity.GalleryImage
	for _, fname := range filenames {
		images = append(images, entity.GalleryImage{
			InvitationID: inv.ID,
			Filename:     fname,
		})
	}
	return s.repo.CreateGallery(images)
}

func (s *invitationService) DeleteGalleryImage(id string) error {
	img, err := s.repo.FindGalleryImageByID(id)
	if err != nil {
		return err
	}

	if img.Filename != "" {
		path := filepath.Join("public", "uploads", img.Filename)
		_ = os.Remove(path) 
	}

	return s.repo.DeleteGalleryImage(id)
}