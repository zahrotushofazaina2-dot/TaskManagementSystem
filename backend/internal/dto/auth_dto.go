package dto

import "time"

// RegisterRequest adalah data yang diterima saat registrasi.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// LoginRequest adalah data yang diterima saat login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// UserResponse adalah data pengguna yang boleh ditampilkan.
type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthResponse adalah respons setelah register atau login.
type AuthResponse struct {
	Message string       `json:"message"`
	Token   string       `json:"token"`
	User    UserResponse `json:"user"`
}
