package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomType struct {
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Description   *string        `gorm:"type:text" json:"description,omitempty"`
	PricePerMonth float64        `gorm:"type:decimal(12,2);not null" json:"price_per_month"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	Rooms      []Room     `gorm:"foreignKey:RoomTypeID" json:"rooms,omitempty"`
	Facilities []Facility `gorm:"many2many:room_type_facilities;" json:"facilities,omitempty"`
}

type Room struct {
	ID         uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	RoomTypeID uuid.UUID      `gorm:"type:char(36);not null;index" json:"room_type_id"`
	RoomNumber string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"room_number"`
	Status     string         `gorm:"type:enum('available','occupied','maintenance');not null;default:'available'" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	RoomType *RoomType `gorm:"foreignKey:RoomTypeID" json:"room_type,omitempty"`
}
