package handler

import (
	"artela-service/internal/entity"
	"artela-service/internal/repository" // Import repo
	"artela-service/internal/service"
	"artela-service/internal/utils" // Import utils response
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (h *InvitationHandler) UploadGallery(c *fiber.Ctx) error {
    slug := c.Params("slug")

    // 1. Parse Multipart Form
    form, err := c.MultipartForm()
    if err != nil {
        return utils.BuildResponse(c, h.errorRepo, "ART-98-001", nil) // Invalid Data
    }

    files := form.File["photos"] // Key-nya 'photos'

    // 2. Validasi Jumlah File (Max 5)
    if len(files) > 5 {
        return utils.BuildResponse(c, h.errorRepo, "ART-98-005", nil)
    }

    var savedUrls []string

    // 3. Validasi Ekstensi & Simpan File
    for _, file := range files {
        // Ambil ekstensi dan lowercase biar validasi aman
        ext := strings.ToLower(filepath.Ext(file.Filename))

        // Validasi: Hanya JPG dan PNG
        if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
            return utils.BuildResponse(c, h.errorRepo, "ART-98-006", nil)
        }

        // Generate nama unik: uuid + extension
        uniqueName := uuid.New().String() + ext
        
        // Tentukan path penyimpanan (folder 'public/uploads' harus dibuat dulu)
        savePath := fmt.Sprintf("./public/uploads/%s", uniqueName)
        
        // Simpan file fisik
        if err := c.SaveFile(file, savePath); err != nil {
            return utils.BuildResponse(c, h.errorRepo, "ART-99-999", nil)
        }

        // Generate URL Public (asumsi domain localhost/server)
        // Nanti di production bisa diganti domain asli
        fullUrl := fmt.Sprintf("/uploads/%s", uniqueName)
        savedUrls = append(savedUrls, fullUrl)
    }

    // 4. Panggil Service untuk simpan URL ke Database
    if err := h.service.AddGalleryImages(slug, savedUrls); err != nil {
        return utils.BuildResponse(c, h.errorRepo, "ART-99-002", nil)
    }

    return utils.BuildResponse(c, h.errorRepo, "ART-00-001", fiber.Map{
        "uploaded_count": len(savedUrls),
        "urls":           savedUrls,
    })
}