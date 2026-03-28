package config

import (
	"log"
	"os"
	"sadewashub-go/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectionDatabase() {
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	database.SetupJoinTable(&models.User{}, "Roles", &models.ModelHasRole{})
	database.SetupJoinTable(&models.Role{}, "Permissions", &models.RoleHasPermission{})

	err = database.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Role{},
		&models.Permission{},
		&models.RoomType{},
		&models.Room{},
		&models.Facility{},
		&models.Promo{},
		&models.Booking{},
		&models.Payment{},
		&models.Ticket{},
		&models.Expense{},
		&models.Attachment{},
		&models.NotificationLog{},
	)

	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	DB = database
	log.Println("Berhasil terhubung ke MYSQL!")
}
