package seed

import (
	"log"

	"my-portfolio/internal/config"
	"my-portfolio/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// seedAdmin creates the default admin user if no admins exist. Safe to call on every startup.
func seedAdmin(db *gorm.DB, cfg config.TypeMyPortfolio) {
	var count int64
	db.Model(&model.Admin{}).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.DefaultPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash default admin password: %v", err)
	}

	admin := model.Admin{
		Username:     cfg.Admin.DefaultUsername,
		Email:        cfg.Admin.DefaultEmail,
		PasswordHash: string(hash),
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}
	log.Printf("Seeded default admin user: %s", admin.Username)
}
