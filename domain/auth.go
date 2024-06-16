package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

const DefaultRole = "user"

type IAuthRepository interface {
	Register(ctx context.Context, authInfo *domain.User) (uuid.UUID, error)
	GetByUsername(ctx context.Context, username string) (*domain.UserAuth, error)
}
