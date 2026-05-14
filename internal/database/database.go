package database

import (
	"log"

	"school-examination/internal/config"
	"school-examination/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	dsn := config.AppConfig.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
	DB = db
	return db
}

// SeedSuperAdmin membuat super admin default jika belum ada
func SeedSuperAdmin(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Where("role = ?", model.RoleSuperAdmin).Count(&count)
	if count > 0 {
		return
	}

	// bcrypt hash dari "password"
	superAdmin := model.User{
		Name:     "Super Admin",
		Email:    "superadmin@school.com",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
		Role:     model.RoleSuperAdmin,
		IsActive: true,
	}

	if err := db.Create(&superAdmin).Error; err != nil {
		log.Printf("[seed] Failed to seed super admin: %v", err)
		return
	}
	log.Printf("[seed] Super admin created: superadmin@school.com / password (id: %s)", superAdmin.ID)
}

// AutoMigrate alias untuk backward compatibility
func AutoMigrate(db *gorm.DB) {
	Migrate(db)
}
