package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.IUserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	dbModel := rDomain.ToUserDB(user)
	dbModel.ID = uuid.New()
	err := r.db.WithContext(ctx).
		Create(dbModel).Error
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}
	return nil
}

func (r *userRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user rDomain.User
	err := r.db.WithContext(ctx).
		First(&user, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting user by id: %w", err)
	}
	return rDomain.ToUserBL(&user), nil
}

func (r *userRepository) GetAll(ctx context.Context, page int) ([]*domain.User, error) {
	var users []*rDomain.User
	err := r.db.WithContext(ctx).
		Limit(PageSize).
		Offset(PageSize * (page - 1)).
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("getting users: %w", err)
	}

	resUsers := make([]*domain.User, 0)
	for _, user := range users {
		resUsers = append(resUsers, rDomain.ToUserBL(user))
	}

	return resUsers, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	dbUser := rDomain.ToUserDB(user)
	err := r.db.WithContext(ctx).
		Save(dbUser).Error

	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	return nil
}

func (r *userRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&rDomain.User{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting user by id: %w", err)
	}
	return nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var dbUser *rDomain.User
	err := r.db.WithContext(ctx).
		Where("login = ?", username).First(&dbUser).Error
	if err != nil {
		return nil, fmt.Errorf("getting user by username: %w", err)
	}
	return rDomain.ToUserBL(dbUser), nil
}
