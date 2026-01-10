package handler

import (
	"artela-service/internal/entity"
	"artela-service/internal/repository"
	"artela-service/internal/service"
	"artela-service/internal/utils"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InvitationHandler struct {
	service   service.InvitationService
	errorRepo repository.ErrorRepository
}

func NewInvitationHandler(service service.InvitationService, errRepo repository.ErrorRepository) *InvitationHandler {
	return &InvitationHandler{
		service:   service,
		errorRepo: errRepo,
	}
}

func (h *InvitationHandler) GetInvitation(c *fiber.Ctx) error {
	slug := c.Params("slug")
	response, err := h.service.GetInvitation(slug)
	if err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-004", nil)
	}
	return utils.BuildResponse(c, h.errorRepo, "ART-00-000", response)
}

func (h *InvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	var input entity.Invitation
	if err := c.BodyParser(&input); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-001", nil)
	}

	if err := h.service.CreateInvitation(&input); err != nil {
		// Cek error duplicate slug (biasanya error dari GORM mengandung string duplicate key)
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}

	return utils.BuildResponse(c, h.errorRepo, "ART-00-001", input)
}

func (h *InvitationHandler) UploadGallery(c *fiber.Ctx) error {
	slug := c.Params("slug")

	form, err := c.MultipartForm()
	if err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-001", nil)
	}

	files := form.File["photos"]
	if len(files) > 5 {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-005", nil)
	}

	var savedUrls []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return utils.BuildResponse(c, h.errorRepo, "ART-98-006", nil)
		}

		uniqueName := uuid.New().String() + ext
		savePath := fmt.Sprintf("./public/uploads/%s", uniqueName)

		if err := c.SaveFile(file, savePath); err != nil {
			return utils.BuildResponse(c, h.errorRepo, "ART-99-999", nil)
		}
		savedUrls = append(savedUrls, fmt.Sprintf("/uploads/%s", uniqueName))
	}

	if err := h.service.AddGalleryImages(slug, savedUrls); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}

	return utils.BuildResponse(c, h.errorRepo, "ART-00-001", fiber.Map{
		"uploaded_count": len(savedUrls),
		"urls":           savedUrls,
	})
}

// Update Invitation (PUT)
func (h *InvitationHandler) UpdateInvitation(c *fiber.Ctx) error {
	slug := c.Params("slug")
	var input entity.Invitation

	if err := c.BodyParser(&input); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-98-001", nil)
	}

	if err := h.service.UpdateInvitation(slug, &input); err != nil {
		if err.Error() == "data tidak ditemukan" {
			return utils.BuildResponse(c, h.errorRepo, "ART-98-004", nil)
		}
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}

	return utils.BuildResponse(c, h.errorRepo, "ART-00-002", nil)
}

// Delete Invitation (DELETE)
func (h *InvitationHandler) DeleteInvitation(c *fiber.Ctx) error {
	slug := c.Params("slug")

	if err := h.service.DeleteInvitation(slug); err != nil {
		if err.Error() == "data tidak ditemukan" {
			return utils.BuildResponse(c, h.errorRepo, "ART-98-004", nil)
		}
		return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
	}

	return utils.BuildResponse(c, h.errorRepo, "ART-00-003", nil)
}