package seeders

import (
	"log"
	"sadewashub-go/internal/config"
	"sadewashub-go/internal/models"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func RunSeeder() {
	permissions := []string{
		"manage-rooms", "view-rooms", "create-booking", "manage-all-bookings",
		"view-own-booking", "verify-ktp", "manage-users", "edit-profile",
		"manage-payments", "view-own-payments", "manage-settings",
	}

	for _, permName := range permissions {
		var perm models.Permission
		config.DB.FirstOrCreate(&perm, models.Permission{
			ID:        uuid.New(),
			Name:      permName,
			GuardName: "api",
		})
	}

	var adminRole models.Role
	config.DB.FirstOrCreate(&adminRole, models.Role{ID: uuid.New(), Name: "admin", GuardName: "api"})

	var tenantRole models.Role
	config.DB.FirstOrCreate(&tenantRole, models.Role{ID: uuid.New(), Name: "tenant", GuardName: "api"})

	var allPerms []models.Permission
	config.DB.Find(&allPerms)
	config.DB.Model(&adminRole).Association("Permissions").Replace(&allPerms)

	var tenantPerms []models.Permission
	config.DB.Where("name IN ?", []string{
		"view-rooms", "create-booking", "view-own-booking", "edit-profile", "view-own-payments",
	}).Find(&tenantPerms)
	config.DB.Model(&tenantRole).Association("Permissions").Replace(&tenantPerms)

	var adminUser models.User
	adminEmail := "admin@sadewascoliving.com"

	if err := config.DB.Where("email = ?", adminEmail).First(&adminUser).Error; err != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("@Gd03072003"), bcrypt.DefaultCost)
		passStr := string(hashedPassword)
		now := time.Now()

		adminUser = models.User{
			ID:              uuid.New(),
			Name:            "Super Admin",
			Email:           adminEmail,
			Password:        &passStr,
			EmailVerifiedAt: &now,
		}

		profile := models.UserProfile{
			ID:       uuid.New(),
			UserID:   adminUser.ID,
			FullName: "Super Admin Sadewas",
		}

		config.DB.Create(&adminUser)
		config.DB.Create(&profile)

		config.DB.Model(&adminUser).Association("Roles").Append(&adminRole)

		log.Println("✅ Seeder: Akun Admin default berhasil dibuat! (admin@sadewascoliving.com / password123)")
	}
}
