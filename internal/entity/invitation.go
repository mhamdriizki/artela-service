package entity

import (
	"time"

	"gorm.io/gorm"
)

type Invitation struct {
	gorm.Model
	Slug          string         `gorm:"uniqueIndex" json:"slug"`
	Theme         string         `json:"theme"` // Field Baru
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
	GalleryImages []GalleryImage `gorm:"foreignKey:InvitationID" json:"gallery_images"`
	Guestbooks    []Guestbook    `gorm:"foreignKey:InvitationID" json:"guestbooks"`
	RSVPs         []RSVP         `gorm:"foreignKey:InvitationID" json:"rsvps"`
}

type GalleryImage struct {
	gorm.Model
	InvitationID uint   `json:"-"`
	Url          string `json:"url"`
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

// --- DTO Response ---

type InvitationResponse struct {
	Slug          string       `json:"slug"`
	Theme         string       `json:"theme"` // Field Baru
	CoupleName    string       `json:"couple_name"`
	GroomName     string       `json:"groom_name"`
	GroomPhoto    string       `json:"groom_photo_url"`
	BrideName     string       `json:"bride_name"`
	BridePhoto    string       `json:"bride_photo_url"`
	YoutubeUrl    string       `json:"youtube_url"`
	Gallery       []string     `json:"gallery"`
	EventDetails  EventDetails `json:"event_details"`
}

type EventDetails struct {
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	Address  string    `json:"address"`
	MapUrl   string    `json:"map_url"`
}