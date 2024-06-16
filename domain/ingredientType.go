package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

type IngredientType struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
}

func (IngredientType) TableName() string {
	return "ingredientType"
}

func ToIngredientTypeDB(ingredientType *domain.IngredientType) *IngredientType {
	return &IngredientType{
		ID:          ingredientType.ID,
		Name:        ingredientType.Name,
		Description: ingredientType.Description,
	}
}

func ToIngredientTypeBL(ingredientType *IngredientType) *domain.IngredientType {
	return &domain.IngredientType{
		ID:          ingredientType.ID,
		Name:        ingredientType.Name,
		Description: ingredientType.Description,
	}
}

type IIngredientTypeRepository interface {
	Create(ctx context.Context, ingredientType *domain.IngredientType) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.IngredientType, error)
	GetAll(ctx context.Context) ([]*domain.IngredientType, error)
	Update(ctx context.Context, measurement *domain.IngredientType) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
