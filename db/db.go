package db

import (
	"log"

	"auth-service/config"
	"auth-service/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	d, err := gorm.Open(sqlite.Open(config.SQLitePath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	DB = d

	// Automigrate models
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	// Seed admin if not exists
	seedAdmin()
}

func seedAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		admin := models.User{
			Email:    "admin@example.com",
			FullName: "Admin",
			Role:     "admin",
		}
		// password: admin123 (hashed)
		if err := admin.SetPassword("admin123"); err != nil {
			log.Printf("failed to create admin seed: %v", err)
			return
		}
		if err := DB.Create(&admin).Error; err != nil {
			log.Printf("failed to insert admin seed: %v", err)
			return
		}
		log.Println("✅ Admin seeded: admin@example.com / admin123")
	}
}
