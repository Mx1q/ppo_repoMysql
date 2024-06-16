package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) rDomain.IAuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Register(ctx context.Context, authInfo *domain.User) (uuid.UUID, error) {
	dbUser := rDomain.ToUserDB(authInfo)
	dbUser.ID = uuid.New()
	res := r.db.WithContext(ctx).
		Create(&dbUser)
	err := res.Error

	if err != nil {
		return uuid.Nil, fmt.Errorf("user registration: %w", err)
	}
	return authInfo.ID, nil
}

func (r *authRepository) GetByUsername(ctx context.Context, username string) (*domain.UserAuth, error) {
	//var dbUser domain.UserAuth
	var dbU rDomain.User

	err := r.db.WithContext(ctx).
		Where("login = ?", username).
		First(&dbU).Error
	if err != nil {
		return nil, fmt.Errorf("getting user by username: %w", err)
	}

	data := new(domain.UserAuth)
	data.ID = dbU.ID
	data.HashedPass = dbU.Password
	data.Role = dbU.Role
	return data, nil
}
