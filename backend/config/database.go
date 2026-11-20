package config

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// ConnectDatabase membuat koneksi ke database MySQL/MariaDB.
func ConnectDatabase() (*sql.DB, error) {
	// Membaca konfigurasi dari file .env.
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("gagal membaca file .env: %w", err)
	}

	// Konfigurasi koneksi database.
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 net.JoinHostPort(os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		DBName:               os.Getenv("DB_NAME"),
		ParseTime:            true,
		Loc:                  time.Local,
		AllowNativePasswords: true,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}

	// Memastikan konfigurasi wajib tidak kosong.
	if cfg.User == "" || cfg.Addr == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("konfigurasi database belum lengkap")
	}

	// Membuka koneksi database.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	// Mengatur connection pool.
	maxOpen, err := getEnvInt("DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		db.Close()
		return nil, err
	}

	maxIdle, err := getEnvInt("DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		db.Close()
		return nil, err
	}

	if maxOpen < 1 || maxIdle < 0 || maxIdle > maxOpen {
		db.Close()
		return nil, fmt.Errorf("konfigurasi connection pool tidak valid")
	}

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)

	// Mengatur durasi maksimal koneksi.
	lifetime := os.Getenv("DB_CONN_MAX_LIFETIME")
	if lifetime == "" {
		lifetime = "5m"
	}

	maxLifetime, err := time.ParseDuration(lifetime)
	if err != nil || maxLifetime <= 0 {
		db.Close()
		return nil, fmt.Errorf("DB_CONN_MAX_LIFETIME tidak valid")
	}

	db.SetConnMaxLifetime(maxLifetime)

	// Memeriksa koneksi database.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("koneksi database gagal: %w", err)
	}

	return db, nil
}

// getEnvInt membaca environment variable dalam bentuk integer.
func getEnvInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s harus berupa angka", key)
	}

	return result, nil
}
