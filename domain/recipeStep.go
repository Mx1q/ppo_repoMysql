package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

type RecipeStep struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	RecipeID    uuid.UUID `gorm:"column:recipeId"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	StepNum     int       `gorm:"column:stepNum"`
}

func (RecipeStep) TableName() string {
	return "recipeStep"
}

func ToStepDB(step *domain.RecipeStep) *RecipeStep {
	return &RecipeStep{
		ID:          step.ID,
		RecipeID:    step.RecipeID,
		Name:        step.Name,
		Description: step.Description,
		StepNum:     step.StepNum,
	}
}

func ToStepBL(step *RecipeStep) *domain.RecipeStep {
	return &domain.RecipeStep{
		ID:          step.ID,
		RecipeID:    step.RecipeID,
		Name:        step.Name,
		Description: step.Description,
		StepNum:     step.StepNum,
	}
}

type IRecipeStepRepository interface {
	Create(ctx context.Context, recipeStep *domain.RecipeStep) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.RecipeStep, error)
	GetAllByRecipeID(ctx context.Context, recipeId uuid.UUID) ([]*domain.RecipeStep, error)
	Update(ctx context.Context, recipeStep *domain.RecipeStep) error
	DeleteById(ctx context.Context, id uuid.UUID) error
	DeleteAllByRecipeID(ctx context.Context, recipeId uuid.UUID) error
}
