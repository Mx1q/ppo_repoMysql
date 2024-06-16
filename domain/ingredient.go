package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

type Ingredient struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	TypeID   uuid.UUID `gorm:"column:type"`
	Name     string    `gorm:"column:name"`
	Calories int       `gorm:"column:calories"`
}

type IngredientLink struct {
	ID           uuid.UUID `gorm:"primaryKey"`
	RecipeId     uuid.UUID `gorm:"column:recipeId"`
	IngredientId uuid.UUID `gorm:"column:ingredientId"`
	Measurement  uuid.UUID `gorm:"column:measurement"`
	amount       int       `gorm:"column:amount"`
}

func (Ingredient) TableName() string {
	return "ingredient"
}

func (IngredientLink) TableName() string {
	return "recipeIngredient"
}

func ToIngredientDB(ingredient *domain.Ingredient) *Ingredient {
	return &Ingredient{
		ID:       ingredient.ID,
		TypeID:   ingredient.TypeID,
		Name:     ingredient.Name,
		Calories: ingredient.Calories,
	}
}

func ToIngredientBL(ingredient *Ingredient) *domain.Ingredient {
	return &domain.Ingredient{
		ID:       ingredient.ID,
		TypeID:   ingredient.TypeID,
		Name:     ingredient.Name,
		Calories: ingredient.Calories,
	}
}

type IIngredientRepository interface {
	Create(ctx context.Context, ingredient *domain.Ingredient) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.Ingredient, error)
	GetAll(ctx context.Context, page int) ([]*domain.Ingredient, int, error)
	GetAllByRecipeId(ctx context.Context, id uuid.UUID) ([]*domain.Ingredient, error)
	Link(ctx context.Context, recipeId uuid.UUID, ingredientId uuid.UUID) (uuid.UUID, error)
	Unlink(ctx context.Context, recipeId uuid.UUID, ingredientId uuid.UUID) error
	Update(ctx context.Context, ingredient *domain.Ingredient) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
