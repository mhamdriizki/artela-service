package entity

import (
	"time"

	"gorm.io/gorm"
)

// Main Entity
type Invitation struct {
	gorm.Model
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
	
	// Perubahan: Rename dari GalleryImages -> Gallery (agar match dengan Preload("Gallery"))
	Gallery       []GalleryImage `gorm:"foreignKey:InvitationID" json:"gallery"` 
	
	Guestbooks    []Guestbook    `gorm:"foreignKey:InvitationID" json:"guestbooks"`
	RSVPs         []RSVP         `gorm:"foreignKey:InvitationID" json:"rsvps"`
}

// Entity Gallery
type GalleryImage struct {
	gorm.Model
	InvitationID uint   `json:"-"`
	// Perubahan: Rename dari Url -> Filename (agar match dengan Service logic)
	Filename     string `json:"filename"` 
}

type Guestbook struct {
	gorm.Model
	InvitationID uint   `json:"invitation_id"`
	Name         string `json:"name"`
	Message      string `json:"message"`
}

type RSVP struct {
	gorm.Model
	InvitationID uint   `json:"invitation_id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Pax          int    `json:"pax"`
}

// --- DTO Responses (Untuk List Dashboard) ---

// Wrapper { "data": [...] }
type InvitationListWrapper struct {
	Data []InvitationListResponse `json:"data"`
}

// Item Response Ringan
type InvitationListResponse struct {
	Slug       string `json:"slug"`
	CoupleName string `json:"couple_name"`
	Theme      string `json:"theme"`
	CreatedAt  string `json:"created_at"`
}

// --- DTO Public (Optional, jika dipakai di response public) ---

type InvitationResponse struct {
	Slug          string       `json:"slug"`
	Theme         string       `json:"theme"`
	CoupleName    string       `json:"couple_name"`
	GroomName     string       `json:"groom_name"`
	GroomPhoto    string       `json:"groom_photo_url"`
	BrideName     string       `json:"bride_name"`
	BridePhoto    string       `json:"bride_photo_url"`
	YoutubeUrl    string       `json:"youtube_url"`
	Gallery       []string     `json:"gallery"` // Array filename string
	EventDetails  EventDetails `json:"event_details"`
}

type EventDetails struct {
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	Address  string    `json:"address"`
	MapUrl   string    `json:"map_url"`
}