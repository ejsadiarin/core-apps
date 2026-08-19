package user

import (
	"context"

	usertypes "github.com/ejsadiarin/coregateway/internal/domain/user/v1"
	"github.com/google/uuid"
)

func GetUserByEmail(ctx context.Context, user *usertypes.User) (*usertypes.User, error) {
	return nil, nil
}

func GetUserByID(ctx context.Context, id uuid.UUID) (*usertypes.User, error) {
	return nil, nil
}

func CreateUser(ctx context.Context, email, passwordHash, role string) (*usertypes.User, error) {
	return nil, nil
}
