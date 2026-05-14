package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"school-examination/config"
	"school-examination/database"
)

func main() {
	// Flags
	fresh  := flag.Bool("fresh", false, "Drop semua tabel lalu recreate (HAPUS SEMUA DATA)")
	seed   := flag.Bool("seed", false, "Jalankan seeder setelah migrate")
	drop   := flag.Bool("drop", false, "Drop semua tabel saja tanpa recreate")
	help   := flag.Bool("help", false, "Tampilkan bantuan")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	// Load config
	config.Load()

	// Koneksi database
	db := database.Connect()

	switch {
	case *drop:
		// Hanya drop, tidak recreate
		confirm("DROP semua tabel (semua data akan hilang)")
		database.Drop(db)
		log.Println("Done. Semua tabel telah dihapus.")

	case *fresh:
		// Drop + recreate + optional seed
		confirm("FRESH migration (semua data akan hilang, tabel akan dibuat ulang)")
		database.Fresh(db)
		if *seed {
			database.Seed(db)
		}
		log.Println("Done. Database berhasil direset.")

	default:
		// AutoMigrate biasa (tambah kolom baru, tidak hapus data)
		log.Println("Running migration (safe mode - data tidak dihapus)...")
		database.Migrate(db)
		if *seed {
			database.Seed(db)
		}
		log.Println("Done.")
	}
}

func confirm(action string) {
	fmt.Printf("\n⚠️  WARNING: %s!\n", action)
	fmt.Print("Ketik 'yes' untuk melanjutkan: ")

	var input string
	fmt.Scanln(&input)

	if input != "yes" {
		fmt.Println("Dibatalkan.")
		os.Exit(0)
	}
}

func printHelp() {
	fmt.Println(`
School Examination - Database Migration Tool

Usage:
  go run cmd/migrate/main.go [flags]

Flags:
  (tanpa flag)       AutoMigrate - tambah kolom baru, data aman
  --fresh            Drop semua tabel + recreate dari awal (DATA HILANG)
  --drop             Drop semua tabel saja
  --seed             Jalankan seeder setelah migrate
  --help             Tampilkan bantuan ini

Contoh:
  go run cmd/migrate/main.go                   # migrate aman
  go run cmd/migrate/main.go --fresh           # reset database
  go run cmd/migrate/main.go --fresh --seed    # reset + isi data awal
  go run cmd/migrate/main.go --seed            # migrate + seed saja
  go run cmd/migrate/main.go --drop            # hapus semua tabel
`)
}
