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
	DeleteGalleryImage(id uint) error // <-- BARU
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

	// Mapping ke list response
	var list []entity.InvitationListResponse
	for _, item := range data {
		list = append(list, entity.InvitationListResponse{
			Slug:       item.Slug,
			CoupleName: item.CoupleName,
			Theme:      item.Theme,
			CreatedAt:  item.CreatedAt.Format("2006-01-02"),
		})
	}
	
	// Init empty slice jika null agar JSON tetap []
	if list == nil {
		list = []entity.InvitationListResponse{}
	}

	// Return dengan Wrapper
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
	// Pastikan ID tetap sama saat update
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

// --- IMPLEMENTASI BARU (DELETE GALLERY) ---

func (s *invitationService) DeleteGalleryImage(id uint) error {
	// 1. Ambil data gambar untuk tahu nama filenya
	img, err := s.repo.FindGalleryImageByID(id)
	if err != nil {
		return err
	}

	// 2. Hapus File Fisik di folder public/uploads
	if img.Filename != "" {
		// Pastikan path sesuai struktur foldermu
		path := filepath.Join("public", "uploads", img.Filename)
		// Hapus file, abaikan error jika file sudah tidak ada (biar DB tetap bersih)
		_ = os.Remove(path) 
	}

	// 3. Hapus Record Database
	return s.repo.DeleteGalleryImage(id)
}