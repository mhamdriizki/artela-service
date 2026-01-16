package entity

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	BaseEntity
	Slug          string    `gorm:"uniqueIndex" json:"slug"`
	Theme         string    `json:"theme"`
	CoupleName    string    `json:"couple_name"`
	
	// Data Mempelai
	GroomName     string    `json:"groom_name"`
	GroomPhoto    string    `json:"groom_photo_url"`
	BrideName     string    `json:"bride_name"`
	BridePhoto    string    `json:"bride_photo_url"`
	
	// Field Baru: Tanggal & Lokasi
	WeddingDate        time.Time `json:"wedding_date"`         // Tanggal Pernikahan Utama
	
	AkadLocation       string    `json:"akad_location"`        // Lokasi Akad
	AkadMapUrl         string    `json:"akad_map_url"`         // Google Maps Akad
	
	ReceptionLocation  string    `json:"reception_location"`   // Lokasi Resepsi
	ReceptionMapUrl    string    `json:"reception_map_url"`    // Google Maps Resepsi
	
	// Multimedia
	YoutubeUrl         string    `json:"youtube_url"`          // Video Utama (jika ada)
	BackgroundMusicUrl string    `json:"background_music_url"` // Backsound (Link Youtube)

	// Relasi
	Gallery       []GalleryImage `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"gallery"`
	Guestbooks    []Guestbook    `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"guestbooks"`
	RSVPs         []RSVP         `gorm:"foreignKey:InvitationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"rsvps"`
}

type GalleryImage struct {
	BaseEntity
	InvitationID uuid.UUID `gorm:"type:char(36)" json:"invitation_id"`
	Filename     string    `json:"filename"`
}

type Guestbook struct {
	BaseEntity
	InvitationID uuid.UUID `gorm:"type:char(36)" json:"invitation_id"`
	Name         string    `json:"name"`
	Message      string    `json:"message"`
}

type RSVP struct {
	BaseEntity
	InvitationID uuid.UUID `gorm:"type:char(36)" json:"invitation_id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	Pax          int       `json:"pax"`
}

// --- DTO ---

type InvitationListWrapper struct {
	Data []InvitationListResponse `json:"data"`
}

type InvitationListResponse struct {
	Slug       string `json:"slug"`
	CoupleName string `json:"couple_name"`
	Theme      string `json:"theme"`
	WeddingDate string `json:"wedding_date"`
	CreatedAt  string `json:"created_at"`
}