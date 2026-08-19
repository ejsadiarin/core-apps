package user

import (
	"context"

	usertypes "github.com/ejsadiarin/coregateway/internal/domain/user/v1"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, user *usertypes.User) (*usertypes.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*usertypes.User, error)
	CreateUser(ctx context.Context, email, passwordHash, role string) (*usertypes.User, error)
}

// User management requests (admin)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	Role     string `json:"role" validate:"required,oneof=guest user admin"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8,max=128"`
	Role     *string `json:"role,omitempty" validate:"omitempty,oneof=guest user admin"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt string    `json:"created_at"`
}
