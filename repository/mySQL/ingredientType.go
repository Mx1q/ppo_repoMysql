package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ingredientTypeRepository struct {
	db *gorm.DB
}

func NewIngredientTypeRepository(db *gorm.DB) rDomain.IIngredientTypeRepository {
	return &ingredientTypeRepository{
		db: db,
	}
}

func (r *ingredientTypeRepository) Create(ctx context.Context, ingredientType *domain.IngredientType) error {
	dbIngredientType := rDomain.ToIngredientTypeDB(ingredientType)
	dbIngredientType.ID = uuid.New()
	err := r.db.WithContext(ctx).
		Create(&dbIngredientType).Error
	if err != nil {
		return fmt.Errorf("creating ingredient type: %w", err)
	}
	return nil
}

func (r *ingredientTypeRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.IngredientType, error) {
	var ingredientType rDomain.IngredientType
	err := r.db.WithContext(ctx).
		First(&ingredientType, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting ingredient type by id: %w", err)
	}
	return rDomain.ToIngredientTypeBL(&ingredientType), nil
}

func (r *ingredientTypeRepository) GetAll(ctx context.Context) ([]*domain.IngredientType, error) {
	var ingredientTypes []*rDomain.IngredientType
	err := r.db.WithContext(ctx).
		Find(&ingredientTypes).Error
	if err != nil {
		return nil, fmt.Errorf("getting salad types: %w", err)
	}

	resIngredientTypes := make([]*domain.IngredientType, 0)
	for _, ingredientType := range ingredientTypes {
		resIngredientTypes = append(resIngredientTypes, rDomain.ToIngredientTypeBL(ingredientType))
	}

	return resIngredientTypes, nil
}

func (r *ingredientTypeRepository) Update(ctx context.Context, measurement *domain.IngredientType) error {
	dbIngredientType := rDomain.ToIngredientTypeDB(measurement)
	err := r.db.WithContext(ctx).
		Save(&dbIngredientType).Error
	if err != nil {
		return fmt.Errorf("updating ingredient type: %w", err)
	}
	return nil
}

func (r *ingredientTypeRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&rDomain.IngredientType{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting ingredient type by id: %w", err)
	}
	return nil
}
