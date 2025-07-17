package models

import "time"

type Property struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Type        string    `json:"type"`
	AddressID   uint      `json:"address_id"`
	OwnerID     uint      `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Area        float64   `json:"area"`
	Rooms       int       `json:"rooms"`
	Bathrooms   int       `json:"bathrooms"`
	Amenities   string    `json:"amenities"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
