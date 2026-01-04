package handler

import (
	"artela-service/internal/entity"
	"artela-service/internal/service"
	"github.com/gofiber/fiber/v2"
)

type InvitationHandler struct {
	service service.InvitationService
}

func NewInvitationHandler(service service.InvitationService) *InvitationHandler {
	return &InvitationHandler{service: service}
}

func (h *InvitationHandler) GetInvitation(c *fiber.Ctx) error {
	slug := c.Params("slug")
	
	response, err := h.service.GetInvitation(slug)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *InvitationHandler) CreateInvitation(c *fiber.Ctx) error {
	var input entity.Invitation
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Data"})
	}

	if err := h.service.CreateInvitation(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data"})
	}

	return c.JSON(fiber.Map{"success": true, "data": input})
}