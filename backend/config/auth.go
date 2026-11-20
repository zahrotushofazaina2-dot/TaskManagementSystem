package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// AuthConfig menyimpan konfigurasi autentikasi.
type AuthConfig struct {
	JWTSecret     string
	JWTExpireHour int
}

// LoadAuthConfig membaca konfigurasi JWT dari environment.
func LoadAuthConfig() (*AuthConfig, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET belum dikonfigurasi")
	}

	hours := 24
	if value := os.Getenv("JWT_EXPIRE_HOURS"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return nil, fmt.Errorf("JWT_EXPIRE_HOURS harus berupa angka positif")
		}
		hours = parsed
	}

	return &AuthConfig{
		JWTSecret:     secret,
		JWTExpireHour: hours,
	}, nil
}

// Expiration mengembalikan durasi masa berlaku token.
func (c *AuthConfig) Expiration() time.Duration {
	return time.Duration(c.JWTExpireHour) * time.Hour
}
