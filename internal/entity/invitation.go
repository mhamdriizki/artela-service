package entity

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	BaseEntity
	Slug          string         `gorm:"uniqueIndex" json:"slug"`
	Theme         string         `json:"theme"`
	CoupleName    string         `json:"couple_name"`
	GroomName     string         `json:"groom_name"`
	GroomPhoto    string         `json:"groom_photo_url"`
	BrideName     string         `json:"bride_name"`
	BridePhoto    string         `json:"bride_photo_url"`
	YoutubeUrl    string         `json:"youtube_url"`
	EventDate     time.Time      `json:"event_date"`
	EventLocation string         `json:"event_location"`
	EventAddress  string         `json:"event_address"`
	MapUrl        string         `json:"map_url"`
	
	// Relation ID tipe UUID
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
	CreatedAt  string `json:"created_at"`
}