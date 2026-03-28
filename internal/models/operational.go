package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:char(36);not null;index" json:"user_id"`
	RoomID      uuid.UUID      `gorm:"type:char(36);not null;index" json:"room_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text;not null" json:"description"`
	Status      string         `gorm:"type:enum('open','in_progress','resolved','rejected');not null;default:'open';index" json:"status"`
	Priority    string         `gorm:"type:enum('low','medium','high','urgent');not null;default:'medium'" json:"priority"`
	PhotoPath   *string        `gorm:"type:varchar(255)" json:"photo_path,omitempty"`
	AdminNote   *string        `gorm:"type:text" json:"admin_note,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Expense struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	Amount      float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	ExpenseDate time.Time      `gorm:"type:date;not null" json:"expense_date"`
	Category    string         `gorm:"type:enum('operational','maintenance','salary','tax','other');not null" json:"category"`
	RoomID      *uuid.UUID     `gorm:"type:char(36);index" json:"room_id,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Model Polimorfik untuk menyimpan file gambar (Bukti TF, Foto Kamar, dll)
type Attachment struct {
	ID             uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	AttachableID   uuid.UUID `gorm:"type:char(36);not null" json:"attachable_id"`
	AttachableType string    `gorm:"type:varchar(255);not null" json:"attachable_type"`
	FilePath       string    `gorm:"type:varchar(255);not null" json:"file_path"`
	FileType       string    `gorm:"type:varchar(255);not null" json:"file_type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type NotificationLog struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	UserID         uuid.UUID  `gorm:"type:char(36);not null;index" json:"user_id"`
	Category       string     `gorm:"type:varchar(255);not null;index" json:"category"`
	Channel        string     `gorm:"type:varchar(255);not null" json:"channel"`
	Title          string     `gorm:"type:varchar(255);not null" json:"title"`
	Message        string     `gorm:"type:text;not null" json:"message"`
	Payload        *string    `gorm:"type:json" json:"payload,omitempty"`
	DeliveryStatus string     `gorm:"type:varchar(255);not null;default:'pending';index" json:"delivery_status"`
	ErrorLog       *string    `gorm:"type:text" json:"error_log,omitempty"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
