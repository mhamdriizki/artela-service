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
	// Response sudah berupa wrapper {data: []} dari service
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
	// Validasi Max 5 Files (Optional, di BE bisa dilepas jika FE sudah handle)
	// if len(files) > 5 { ... }

	var filenames []string
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		// Validasi Extension Sederhana
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			// Skip atau Return Error (Disini kita return error biar strict)
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

// --- IMPLEMENTASI BARU (HANDLER DELETE GALLERY) ---

func (h *InvitationHandler) DeleteGalleryImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-001", nil)
	}

	if err := h.service.DeleteGalleryImage(uint(id)); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}

	return utils.BuildResponse(c, h.errorRepo, "ART-00-003", nil)
}