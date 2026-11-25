package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"

	"taskmanagement/internal/model"
	"taskmanagement/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// seedAdmin membuat akun ADMIN melalui input terminal.
func seedAdmin(db *sql.DB) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("================================")
	fmt.Println("       ADMIN DATABASE SEEDER")
	fmt.Println("================================")

	name, err := readInput(reader, "Nama Admin: ")
	if err != nil {
		return err
	}

	name = strings.TrimSpace(name)

	if len(name) < 2 || len(name) > 100 {
		return fmt.Errorf("nama harus memiliki 2 sampai 100 karakter")
	}

	email, err := readInput(reader, "Email Admin: ")
	if err != nil {
		return err
	}

	email = strings.ToLower(strings.TrimSpace(email))

	if len(email) > 150 {
		return fmt.Errorf("email maksimal 150 karakter")
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || parsedEmail.Address != email {
		return fmt.Errorf("format email tidak valid")
	}

	password, err := readPassword("Password Admin: ")
	if err != nil {
		return err
	}

	if len(password) < 12 || len(password) > 72 {
		return fmt.Errorf("password harus memiliki 12 sampai 72 byte")
	}

	confirmPassword, err := readPassword("Konfirmasi Password: ")
	if err != nil {
		return err
	}

	if password != confirmPassword {
		return fmt.Errorf("konfirmasi password tidak cocok")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	userRepo := repository.NewUserRepository(db)

	existingUser, err := userRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("gagal memeriksa pengguna: %w", err)
	}

	if existingUser != nil {
		if existingUser.Role == "ADMIN" ||
			existingUser.Role == "SUPER_ADMIN" {
			fmt.Println("Akun administrator dengan email tersebut sudah tersedia.")
			return nil
		}

		return fmt.Errorf("email sudah digunakan oleh akun non-Admin")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("gagal melakukan hash password: %w", err)
	}

	admin := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "ADMIN",
		IsActive: true,
	}

	if err := userRepo.Create(ctx, admin); err != nil {
		return fmt.Errorf("gagal membuat Admin: %w", err)
	}

	fmt.Println()
	fmt.Println("================================")
	fmt.Println("   ADMIN BERHASIL DIBUAT")
	fmt.Println("================================")
	fmt.Printf("Admin ID : %d\n", admin.ID)
	fmt.Printf("Nama     : %s\n", admin.Name)
	fmt.Printf("Email    : %s\n", admin.Email)
	fmt.Printf("Role     : %s\n", admin.Role)
	fmt.Printf("Status   : Aktif\n")

	return nil
}

// readInput membaca input teks dari terminal.
func readInput(
	reader *bufio.Reader,
	prompt string,
) (string, error) {
	fmt.Print(prompt)

	value, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("gagal membaca input: %w", err)
	}

	return strings.TrimSpace(value), nil
}

// readPassword membaca password tanpa menampilkannya.
func readPassword(prompt string) (string, error) {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		return "", fmt.Errorf(
			"input password harus dilakukan melalui terminal interaktif",
		)
	}

	fmt.Print(prompt)

	password, err := term.ReadPassword(fd)
	fmt.Println()

	if err != nil {
		return "", fmt.Errorf("gagal membaca password: %w", err)
	}

	return string(password), nil
}
