package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

type Salad struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	AuthorID    uuid.UUID `gorm:"column:authorId"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
}

func (Salad) TableName() string {
	return "salad"
}

func ToSaladDB(salad *domain.Salad) *Salad {
	return &Salad{
		ID:          salad.ID,
		AuthorID:    salad.AuthorID,
		Name:        salad.Name,
		Description: salad.Description,
	}
}

func ToSaladBL(salad *Salad) *domain.Salad {
	return &domain.Salad{
		ID:          salad.ID,
		AuthorID:    salad.AuthorID,
		Name:        salad.Name,
		Description: salad.Description,
	}
}

type ISaladRepository interface {
	Create(ctx context.Context, salad *domain.Salad) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Salad, error)
	GetAll(ctx context.Context, filter *domain.RecipeFilter, page int) ([]*domain.Salad, int, error)
	GetAllByUserId(ctx context.Context, id uuid.UUID) ([]*domain.Salad, error)
	GetAllRatedByUser(ctx context.Context, userId uuid.UUID, page int) ([]*domain.Salad, int, error)
	Update(ctx context.Context, salad *domain.Salad) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
