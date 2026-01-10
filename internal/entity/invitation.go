package entity

import (
	"time"

	"gorm.io/gorm"
)

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
	GalleryImages []GalleryImage `gorm:"foreignKey:InvitationID" json:"gallery_images"`
	Guestbooks    []Guestbook    `gorm:"foreignKey:InvitationID" json:"guestbooks"`
	RSVPs         []RSVP         `gorm:"foreignKey:InvitationID" json:"rsvps"`
}

type GalleryImage struct {
	gorm.Model
	InvitationID uint   `json:"-"` // Hide ID relasi di JSON
	Url          string `json:"url"`
}

// Guestbook Entity
type Guestbook struct {
	gorm.Model
	InvitationID uint   `json:"invitation_id"`
	Name         string `json:"name"`
	Message      string `json:"message"`
}

// RSVP Entity
type RSVP struct {
	gorm.Model
	InvitationID uint   `json:"invitation_id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Pax          int    `json:"pax"`
}

// --- DTO (Data Transfer Object) untuk Response ke Frontend ---

// InvitationResponse (Format JSON Output)
type InvitationResponse struct {
	Slug          string       `json:"slug"`
	Theme         string       `json:"theme"`
	CoupleName    string       `json:"couple_name"`
	GroomName     string       `json:"groom_name"`
	GroomPhoto    string       `json:"groom_photo_url"`
	BrideName     string       `json:"bride_name"`
	BridePhoto    string       `json:"bride_photo_url"`
	YoutubeUrl    string       `json:"youtube_url"`
	Gallery       []string     `json:"gallery"`       // Array string URL simpel
	EventDetails  EventDetails `json:"event_details"` // Nested JSON snake_case
}

type EventDetails struct {
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	Address  string    `json:"address"`
	MapUrl   string    `json:"map_url"`
}