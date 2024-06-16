package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

type Recipe struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	SaladID          uuid.UUID `gorm:"column:saladId"`
	Status           int       `gorm:"column:status"`
	NumberOfServings int       `gorm:"column:numberOfServings"`
	TimeToCook       int       `gorm:"column:timeToCook"`
	Rating           float32   `gorm:"column:rating"`
}

func (Recipe) TableName() string {
	return "recipe"
}

func ToRecipeDB(recipe *domain.Recipe) *Recipe {
	return &Recipe{
		ID:               recipe.ID,
		SaladID:          recipe.SaladID,
		Status:           recipe.Status,
		NumberOfServings: recipe.NumberOfServings,
		TimeToCook:       recipe.TimeToCook,
		Rating:           recipe.Rating,
	}
}

func ToRecipeBL(recipe *Recipe) *domain.Recipe {
	return &domain.Recipe{
		ID:               recipe.ID,
		SaladID:          recipe.SaladID,
		Status:           recipe.Status,
		NumberOfServings: recipe.NumberOfServings,
		TimeToCook:       recipe.TimeToCook,
		Rating:           recipe.Rating,
	}
}

type IRecipeRepository interface {
	Create(ctx context.Context, recipe *domain.Recipe) (uuid.UUID, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.Recipe, error)
	GetBySaladId(ctx context.Context, saladId uuid.UUID) (*domain.Recipe, error)
	GetAll(ctx context.Context, filter *domain.RecipeFilter, page int) ([]*domain.Recipe, error)
	Update(ctx context.Context, recipe *domain.Recipe) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
