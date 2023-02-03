package model

type LoginRequest struct {
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     Role   `json:"role" validate:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     Role   `json:"role" validate:"required"`
}
