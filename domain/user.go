package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"net/mail"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	Name     string    `gorm:"column:name"`
	Username string    `gorm:"column:login"`
	Password string    `gorm:"column:password"`
	Email    string    `gorm:"column:email"`
	Role     string    `gorm:"column:role"`
}

func (User) TableName() string {
	return "user"
}

func ToUserDB(user *domain.User) *User {
	return &User{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email.Address,
		Role:     user.Role,
	}
}

func ToUserBL(user *User) *domain.User {
	return &domain.User{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
		Email:    mail.Address{Address: user.Email},
		Role:     user.Role,
	}
}

type IUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetAll(ctx context.Context, page int) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
