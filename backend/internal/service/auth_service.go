package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"taskmanagement/config"
	"taskmanagement/internal/dto"
	"taskmanagement/internal/model"
	"taskmanagement/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidInput       = errors.New("data tidak valid")
	ErrEmailAlreadyUsed   = errors.New("email sudah terdaftar")
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrAccountInactive    = errors.New("akun tidak aktif")
)

// AuthService menangani logika bisnis autentikasi.
type AuthService struct {
	userRepo *repository.UserRepository
	authCfg  *config.AuthConfig
}

// NewAuthService membuat instance AuthService.
func NewAuthService(
	userRepo *repository.UserRepository,
	authCfg *config.AuthConfig,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		authCfg:  authCfg,
	}
}

// Register mendaftarkan pengguna baru.
func (s *AuthService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) (*dto.AuthResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if len(name) < 2 || len(name) > 100 {
		return nil, ErrInvalidInput
	}

	if len(email) > 150 {
		return nil, ErrInvalidInput
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidInput
	}

	if len(req.Password) < 8 || len(req.Password) > 72 {
		return nil, ErrInvalidInput
	}

	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyUsed
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan hash password: %w", err)
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "USER",
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Message: "Registrasi berhasil",
		User:    toUserResponse(user),
	}, nil
}

// Login memverifikasi kredensial pengguna dan membuat JWT.
func (s *AuthService) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if email == "" || req.Password == "" || len(req.Password) > 72 {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Message: "Login berhasil",
		Token:   token,
		User:    toUserResponse(user),
	}, nil
}

// generateToken membuat JWT dengan identitas pengguna.
func (s *AuthService) generateToken(user *model.User) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"iat":     now.Unix(),
		"exp":     now.Add(s.authCfg.Expiration()).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.authCfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("gagal membuat JWT: %w", err)
	}

	return signedToken, nil
}

// toUserResponse mengubah model menjadi DTO yang aman untuk dikirim.
func toUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Me mengambil data pengguna yang sedang login.
func (s *AuthService) Me(
	ctx context.Context,
	userID int64,
) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil || !user.IsActive {
		return nil, fmt.Errorf("pengguna tidak ditemukan atau tidak aktif")
	}

	response := toUserResponse(user)
	return &response, nil
}
