package main

import (
	"context"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	"taskmanagement/config"
	"taskmanagement/internal/model"
	"taskmanagement/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Membuka koneksi database dan memuat .env.
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	// Membaca konfigurasi Admin.
	name := strings.TrimSpace(os.Getenv("ADMIN_NAME"))
	email := strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_EMAIL")))
	password := os.Getenv("ADMIN_PASSWORD")

	// Validasi konfigurasi.
	if name == "" || email == "" || password == "" {
		log.Fatal("ADMIN_NAME, ADMIN_EMAIL, dan ADMIN_PASSWORD wajib diisi di .env")
	}

	if len(name) < 2 || len(name) > 100 {
		log.Fatal("ADMIN_NAME harus memiliki 2 sampai 100 karakter")
	}

	if len(email) > 150 {
		log.Fatal("ADMIN_EMAIL maksimal 150 karakter")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		log.Fatal("ADMIN_EMAIL tidak valid")
	}

	if len(password) < 12 || len(password) > 72 {
		log.Fatal("ADMIN_PASSWORD harus memiliki 12 sampai 72 karakter")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userRepo := repository.NewUserRepository(db)

	// Memeriksa apakah email sudah terdaftar.
	existingUser, err := userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Fatalf("Gagal memeriksa pengguna: %v", err)
	}

	if existingUser != nil {
		if existingUser.Role == "ADMIN" {
			log.Println("Admin dengan email tersebut sudah tersedia.")
			return
		}

		log.Fatal("Email sudah digunakan oleh akun non-Admin. Seeder dibatalkan.")
	}

	// Melakukan hash password.
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatalf("Gagal melakukan hash password: %v", err)
	}

	// Membuat akun Admin.
	admin := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "ADMIN",
		IsActive: true,
	}

	if err := userRepo.Create(ctx, admin); err != nil {
		log.Fatalf("Gagal membuat Admin: %v", err)
	}

	fmt.Println("Seeder Admin berhasil dijalankan.")
	fmt.Printf("Admin ID: %d\n", admin.ID)
	fmt.Printf("Admin Email: %s\n", admin.Email)
	fmt.Println("Role: ADMIN")
}
