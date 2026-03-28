package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Promo struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Code         string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"code"`
	Type         string    `gorm:"type:enum('percentage','fixed');not null" json:"type"`
	RewardAmount float64   `gorm:"type:decimal(12,2);not null" json:"reward_amount"`
	StartDate    time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate      time.Time `gorm:"type:date;not null" json:"end_date"`
	Limit        *int      `gorm:"type:int" json:"limit,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Booking struct {
	ID             uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID         uuid.UUID      `gorm:"type:char(36);not null;index" json:"user_id"`
	RoomID         uuid.UUID      `gorm:"type:char(36);not null;index" json:"room_id"`
	PromoID        *uuid.UUID     `gorm:"type:char(36);index" json:"promo_id,omitempty"`
	CheckInDate    time.Time      `gorm:"type:date;not null;index" json:"check_in_date"`
	CheckOutDate   *time.Time     `gorm:"type:date;index" json:"check_out_date,omitempty"`
	Status         string         `gorm:"type:enum('pending','confirmed','cancelled','completed');not null;default:'pending';index" json:"status"`
	Reason         *string        `gorm:"type:varchar(255)" json:"reason,omitempty"`
	Notes          *string        `gorm:"type:text" json:"notes,omitempty"`
	TotalAmount    float64        `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	DiscountAmount float64        `gorm:"type:decimal(12,2);not null;default:0.00" json:"discount_amount"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Room    *Room    `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Payment *Payment `gorm:"foreignKey:BookingID" json:"payment,omitempty"`
}

type Payment struct {
	ID            uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	BookingID     uuid.UUID      `gorm:"type:char(36);not null;index" json:"booking_id"`
	ExternalID    *string        `gorm:"type:varchar(255);uniqueIndex" json:"external_id,omitempty"`
	PaymentMethod *string        `gorm:"type:varchar(255)" json:"payment_method,omitempty"`
	Amount        float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Status        string         `gorm:"type:enum('pending','paid','failed','expired');not null;default:'pending';index" json:"status"`
	CheckoutUrl   *string        `gorm:"type:varchar(255)" json:"checkout_url,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
