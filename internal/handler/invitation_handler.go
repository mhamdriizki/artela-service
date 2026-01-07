package handler

import (
	"artela-service/internal/entity"
	"artela-service/internal/repository" // Import repo
	"artela-service/internal/service"
	"artela-service/internal/utils" // Import utils response

	"github.com/gofiber/fiber/v2"
)

type InvitationHandler struct {
	service   service.InvitationService
	errorRepo repository.ErrorRepository // Tambah ini
}

// Update Constructor
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
		// Return Error Response Standard (ART-99-999 atau buat kode khusus misal ART-40-404)
		return utils.BuildResponse(c, h.errorRepo, "ART-99-999", nil)
	}

	// Return Success Response Standard
	return utils.BuildResponse(c, h.errorRepo, "ART-00-000", response)
}

func (h *InvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	var input entity.Invitation
	if err := c.BodyParser(&input); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-999", nil)
	}

	if err := h.service.CreateInvitation(&input); err != nil {
		return utils.BuildResponse(c, h.errorRepo, "ART-99-999", nil)
	}

	return utils.BuildResponse(c, h.errorRepo, "ART-00-000", input)
}