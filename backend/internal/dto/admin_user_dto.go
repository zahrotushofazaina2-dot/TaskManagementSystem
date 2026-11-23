package dto

import "time"

// AdminUserListRequest menampung parameter pencarian dan pagination.
type AdminUserListRequest struct {
	Search   string
	Role     string
	IsActive *bool
	Page     int
	Limit    int
}

// AdminUserResponse adalah data pengguna untuk kebutuhan Admin.
type AdminUserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AdminUserListResponse adalah respons daftar pengguna.
type AdminUserListResponse struct {
	Users      []AdminUserResponse `json:"users"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}

// UpdateUserRoleRequest adalah request perubahan role pengguna.
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=USER ADMIN"`
}

// UpdateUserStatusRequest adalah request perubahan status pengguna.
type UpdateUserStatusRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}
