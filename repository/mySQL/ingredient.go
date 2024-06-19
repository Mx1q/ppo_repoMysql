package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ingredientRepository struct {
	db *gorm.DB
}

func NewIngredientRepository(db *gorm.DB) rDomain.IIngredientRepository {
	return &ingredientRepository{
		db: db,
	}
}

func (r *ingredientRepository) Create(ctx context.Context, ingredient *domain.Ingredient) error {
	dbIngredient := rDomain.ToIngredientDB(ingredient)
	dbIngredient.ID = uuid.New()
	err := r.db.WithContext(ctx).
		Create(&dbIngredient).Error
	if err != nil {
		return fmt.Errorf("creating ingredient: %w", err)
	}
	return nil
}

func (r *ingredientRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Ingredient, error) {
	var ingredient rDomain.Ingredient
	err := r.db.WithContext(ctx).
		First(&ingredient, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting ingredient by id: %w", err)
	}
	return rDomain.ToIngredientBL(&ingredient), nil
}

func (r *ingredientRepository) GetAll(ctx context.Context, page int) ([]*domain.Ingredient, int, error) {
	var ingredients []*rDomain.Ingredient
	err := r.db.WithContext(ctx).
		Limit(PageSize).
		Offset(PageSize * (page - 1)).
		Find(&ingredients).Error
	if err != nil {
		return nil, 0, fmt.Errorf("getting ingredients: %w", err)
	}

	resIngredients := make([]*domain.Ingredient, 0)
	for _, ingredient := range ingredients {
		resIngredients = append(resIngredients, rDomain.ToIngredientBL(ingredient))
	}

	var tmp []*rDomain.Ingredient
	count := r.db.WithContext(ctx).
		Find(&tmp).RowsAffected
	numPages := count / PageSize
	if count%PageSize != 0 {
		numPages++
	}

	return resIngredients, int(numPages), nil
}

func (r *ingredientRepository) GetAllByRecipeId(ctx context.Context, id uuid.UUID) ([]*domain.Ingredient, error) {
	var ingredients []*rDomain.Ingredient
	var ingredientIds []uuid.UUID

	err := r.db.WithContext(ctx).
		Table("recipeIngredient").
		Where("recipeId = ?", id).
		Select("ingredientId").Scan(&ingredientIds).Error
	if err != nil {
		return nil, fmt.Errorf("getting recipe ingredient ids: %w", err)
	}

	if len(ingredientIds) != 0 {
		err = r.db.WithContext(ctx).
			Find(&ingredients, ingredientIds).Error
		if err != nil {
			return nil, fmt.Errorf("getting recipe ingredients: %w", err)
		}
	}

	resIngredients := make([]*domain.Ingredient, 0)
	for _, ingredient := range ingredients {
		resIngredients = append(resIngredients, rDomain.ToIngredientBL(ingredient))
	}

	return resIngredients, nil
}

func (r *ingredientRepository) Update(ctx context.Context, ingredient *domain.Ingredient) error {
	dbIngredient := rDomain.ToIngredientDB(ingredient)
	err := r.db.WithContext(ctx).
		Save(&dbIngredient).Error
	if err != nil {
		return fmt.Errorf("updating ingredient: %w", err)
	}
	return nil
}

func (r *ingredientRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&rDomain.Ingredient{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting ingredient by id: %w", err)
	}
	return nil
}

func (r *ingredientRepository) Link(ctx context.Context, recipeId uuid.UUID, ingredientId uuid.UUID) (uuid.UUID, error) {
	link := rDomain.IngredientLink{
		ID:           uuid.New(), // TODO
		RecipeId:     recipeId,
		IngredientId: ingredientId,
	}
	err := r.db.WithContext(ctx).
		Select("recipeId", "ingredientId").
		Create(&link).Error
	if err != nil {
		return uuid.Nil, fmt.Errorf("linking ingredient: %w", err)
	}
	return link.ID, nil
}

func (r *ingredientRepository) Unlink(ctx context.Context, recipeId uuid.UUID, ingredientId uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("recipeID = ?", recipeId).
		Where("ingredientId = ?", ingredientId).
		Delete(&rDomain.IngredientLink{}).Error
	if err != nil {
		return fmt.Errorf("unlinking ingredient by recipe id: %w", err)
	}
	return nil
}
