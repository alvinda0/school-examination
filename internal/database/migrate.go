package database

import (
	"log"

	"school-examination/internal/model"

	"gorm.io/gorm"
)

// urutan tabel harus dari child ke parent saat drop (hindari FK constraint error)
var allModels = []interface{}{
	&model.StudentAnswer{},
	&model.ExamSubmission{},
	&model.ExamQuestion{},
	&model.Exam{},
	&model.StudentClass{},
	&model.Option{},
	&model.Question{},
	&model.Subject{},
	&model.Class{},
	&model.User{},
}

// Migrate menjalankan AutoMigrate (tambah kolom baru, tidak hapus data)
func Migrate(db *gorm.DB) {
	log.Println("[migrate] Running AutoMigrate...")

	// urutan migrate: parent dulu baru child
	ordered := []interface{}{
		&model.User{},
		&model.Subject{},
		&model.Class{},
		&model.Question{},
		&model.Option{},
		&model.StudentClass{},
		&model.Exam{},
		&model.ExamQuestion{},
		&model.ExamSubmission{},
		&model.StudentAnswer{},
	}

	if err := db.AutoMigrate(ordered...); err != nil {
		log.Fatalf("[migrate] AutoMigrate failed: %v", err)
	}
	log.Println("[migrate] AutoMigrate completed successfully")
}

// Fresh drop semua tabel lalu recreate dari awal (HAPUS SEMUA DATA)
func Fresh(db *gorm.DB) {
	log.Println("[migrate:fresh] WARNING: Dropping all tables...")
	Drop(db)
	log.Println("[migrate:fresh] Recreating tables...")
	Migrate(db)
	log.Println("[migrate:fresh] Fresh migration completed")
}

// Drop menghapus semua tabel (urutan child → parent agar tidak FK error)
func Drop(db *gorm.DB) {
	log.Println("[migrate:drop] Dropping tables...")

	// Disable foreign key checks sementara (PostgreSQL)
	db.Exec("SET session_replication_role = 'replica'")

	for _, m := range allModels {
		if err := db.Migrator().DropTable(m); err != nil {
			log.Printf("[migrate:drop] Warning: could not drop table for %T: %v", m, err)
		}
	}

	// Re-enable foreign key checks
	db.Exec("SET session_replication_role = 'origin'")

	log.Println("[migrate:drop] All tables dropped")
}

// Seed mengisi data awal ke database
func Seed(db *gorm.DB) {
	log.Println("[seed] Running seeders...")
	SeedSuperAdmin(db)
	SeedSampleData(db)
	log.Println("[seed] Seeding completed")
}

// SeedSampleData mengisi data contoh untuk development
func SeedSampleData(db *gorm.DB) {
	// Seed subjects
	subjects := []model.Subject{
		{Name: "Matematika", Code: "MTK", Description: "Mata pelajaran Matematika"},
		{Name: "Bahasa Indonesia", Code: "BIN", Description: "Mata pelajaran Bahasa Indonesia"},
		{Name: "Bahasa Inggris", Code: "BIG", Description: "Mata pelajaran Bahasa Inggris"},
		{Name: "IPA", Code: "IPA", Description: "Ilmu Pengetahuan Alam"},
		{Name: "IPS", Code: "IPS", Description: "Ilmu Pengetahuan Sosial"},
	}

	for i := range subjects {
		var count int64
		db.Model(&model.Subject{}).Where("code = ?", subjects[i].Code).Count(&count)
		if count == 0 {
			db.Create(&subjects[i])
			log.Printf("[seed] Subject created: %s", subjects[i].Name)
		}
	}

	// Seed classes
	classes := []model.Class{
		{Name: "X-A", Grade: "X"},
		{Name: "X-B", Grade: "X"},
		{Name: "XI-A", Grade: "XI"},
		{Name: "XI-B", Grade: "XI"},
		{Name: "XII-A", Grade: "XII"},
	}

	for i := range classes {
		var count int64
		db.Model(&model.Class{}).Where("name = ?", classes[i].Name).Count(&count)
		if count == 0 {
			db.Create(&classes[i])
			log.Printf("[seed] Class created: %s", classes[i].Name)
		}
	}
}
