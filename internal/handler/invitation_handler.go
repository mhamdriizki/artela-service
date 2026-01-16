package handler

import (
	"artela-service/internal/entity"
	"artela-service/internal/repository"
	"artela-service/internal/service"
	"artela-service/internal/utils"
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

type InvitationHandler struct {
	service   service.InvitationService
	errorRepo repository.ErrorRepository
}

func NewInvitationHandler(service service.InvitationService, errorRepo repository.ErrorRepository) *InvitationHandler {
	return &InvitationHandler{service: service, errorRepo: errorRepo}
}

func (h *InvitationHandler) GetAllInvitations(c *fiber.Ctx) error {
	response, err := h.service.GetAllInvitations()
	if err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}
	return utils.BuildResponse(c, h.errorRepo, "ART-00-000", response)
}

func (h *InvitationHandler) GetInvitation(c *fiber.Ctx) error {
	slug := c.Params("slug")
	inv, err := h.service.GetInvitation(slug)
	if err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-004", nil)
	}
	return utils.BuildResponse(c, h.errorRepo, "ART-00-000", inv)
}

func (h *InvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	var req entity.Invitation
	if err := c.BodyParser(&req); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-001", nil)
	}
	if err := h.service.CreateInvitation(&req); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}
	return utils.BuildResponse(c, h.errorRepo, "ART-00-001", req)
}

func (h *InvitationHandler) UpdateInvitation(c *fiber.Ctx) error {
	slug := c.Params("slug")
	var req entity.Invitation
	if err := c.BodyParser(&req); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-001", nil)
	}
	if err := h.service.UpdateInvitation(slug, &req); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}
	return utils.BuildResponse(c, h.errorRepo, "ART-00-002", nil)
}

func (h *InvitationHandler) DeleteInvitation(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if err := h.service.DeleteInvitation(slug); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}
	return utils.BuildResponse(c, h.errorRepo, "ART-00-003", nil)
}

func (h *InvitationHandler) UploadGallery(c *fiber.Ctx) error {
	slug := c.Params("slug")
	form, err := c.MultipartForm()
	if err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-001", nil)
	}

	files := form.File["photos"]
	
	// UPDATE: MAX 7 FILES
	if len(files) > 7 {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-005", nil)
	}

	var filenames []string
	const MaxFileSize = 2 * 1024 * 1024 // 2 MB

	for _, file := range files {
		if file.Size > MaxFileSize {
			return utils.BuildResponse(c, h.errorRepo, "ART-98-007", nil)
		}

		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return utils.BuildResponse(c, h.errorRepo, "ART-98-006", nil)
		}
		
		newFilename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
		if err := c.SaveFile(file, fmt.Sprintf("./public/uploads/%s", newFilename)); err != nil {
			return utils.BuildResponse(c, h.errorRepo, "ART-99-999", nil)
		}
		filenames = append(filenames, newFilename)
	}

	if err := h.service.UploadGallery(slug, filenames); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}

	return utils.BuildResponse(c, h.errorRepo, "ART-00-001", fiber.Map{"uploaded_count": len(filenames)})
}

func (h *InvitationHandler) DeleteGalleryImage(c *fiber.Ctx) error {
	id := c.Params("id") 
	if id == "" {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-001", nil)
	}

	if err := h.service.DeleteGalleryImage(id); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}

	return utils.BuildResponse(c, h.errorRepo, "ART-00-003", nil)
}