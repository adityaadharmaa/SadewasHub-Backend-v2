package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Email           string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	Password        *string        `gorm:"type:varchar(255)" json:"-"` // json:"-" agar tidak ikut terkirim ke frontend
	SocialID        *string        `gorm:"type:varchar(255)" json:"social_id,omitempty"`
	SocialType      *string        `gorm:"type:varchar(255)" json:"social_type,omitempty"`
	RememberToken   *string        `gorm:"type:varchar(100)" json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relasi
	Profile *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`

	Roles []Role `gorm:"many2many:model_has_roles;joinForeignKey:model_id;joinReferences:role_id" json:"roles,omitempty"`
}

type UserProfile struct {
	ID                    uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID                uuid.UUID      `gorm:"type:char(36);not null;index" json:"user_id"`
	FullName              string         `gorm:"type:varchar(255);not null;index" json:"full_name"`
	Nik                   *string        `gorm:"type:varchar(16);uniqueIndex" json:"nik,omitempty"`
	Gender                *string        `gorm:"type:enum('male','female')" json:"gender,omitempty"`
	BirthDate             *time.Time     `gorm:"type:date" json:"birth_date,omitempty"`
	PhoneNumber           *string        `gorm:"type:varchar(255);index" json:"phone_number,omitempty"`
	Address               *string        `gorm:"type:text" json:"address,omitempty"`
	Occupation            *string        `gorm:"type:varchar(255)" json:"occupation,omitempty"`
	EmergencyContactName  *string        `gorm:"type:varchar(255)" json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string        `gorm:"type:varchar(255)" json:"emergency_contact_phone,omitempty"`
	KtpPath               *string        `gorm:"type:varchar(255)" json:"ktp_path,omitempty"`
	IsVerified            bool           `gorm:"type:tinyint(1);not null;default:0" json:"is_verified"`
	AdminNote             *string        `gorm:"type:text" json:"admin_note,omitempty"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// Spatie Roles & Permissions (Disederhanakan)
type Role struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	GuardName string    `gorm:"type:varchar(255);not null" json:"guard_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Permissions []Permission `gorm:"many2many:role_has_permissions;joinForeignKey:role_id;joinReferences:permission_id" json:"permissions,omitempty"`
}

type Permission struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	GuardName string    `gorm:"type:varchar(255);not null" json:"guard_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ModelHasRole struct {
	// HAPUS tulisan ;primaryKey di ketiga baris ini 👇
	RoleID    uuid.UUID `gorm:"type:char(36);column:role_id"`
	ModelType string    `gorm:"type:varchar(255);column:model_type"`
	ModelID   uuid.UUID `gorm:"type:char(36);column:model_id"`
}

func (ModelHasRole) TableName() string {
	return "model_has_roles"
}

func (m *ModelHasRole) BeforeSave(tx *gorm.DB) (err error) {
	if m.ModelType == "" {
		m.ModelType = "App\\Models\\User"
	}
	return
}

type RoleHasPermission struct {
	// HAPUS tulisan ;primaryKey di kedua baris ini 👇
	PermissionID uuid.UUID `gorm:"type:char(36);column:permission_id"`
	RoleID       uuid.UUID `gorm:"type:char(36);column:role_id"`
}

func (RoleHasPermission) TableName() string {
	return "role_has_permissions"
}
