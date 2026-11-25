
package main

import (
	"log"

	"taskmanagement/config"
)

func main() {
	// Menghubungkan aplikasi ke database.
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Menjalankan Admin Seeder.
	if err := seedAdmin(db); err != nil {
		log.Fatalf("Admin seeder failed: %v", err)
	}

	log.Println("Admin seeder completed.")
}
