package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type recipeRepository struct {
	db *gorm.DB
}

func NewRecipeRepository(db *gorm.DB) rDomain.IRecipeRepository {
	return &recipeRepository{
		db: db,
	}
}

func (r *recipeRepository) Create(ctx context.Context, recipe *domain.Recipe) (uuid.UUID, error) {
	dbRecipe := rDomain.ToRecipeDB(recipe)
	dbRecipe.ID = uuid.New()
	err := r.db.WithContext(ctx).
		Create(&dbRecipe).Error
	if err != nil {
		return uuid.Nil, fmt.Errorf("creating recipe: %w", err)
	}
	return dbRecipe.ID, nil
}

func (r *recipeRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
	var recipe rDomain.Recipe
	err := r.db.WithContext(ctx).
		First(&recipe, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting recipe by id: %w", err)
	}
	return rDomain.ToRecipeBL(&recipe), nil
}

func (r *recipeRepository) GetBySaladId(ctx context.Context, saladId uuid.UUID) (*domain.Recipe, error) {
	var recipe rDomain.Recipe
	err := r.db.WithContext(ctx).
		Where("saladId = ?", saladId).
		First(&recipe).Error
	if err != nil {
		return nil, fmt.Errorf("getting recipe by salad id: %w", err)
	}
	return rDomain.ToRecipeBL(&recipe), nil
}

func (r *recipeRepository) GetAll(ctx context.Context, filter *domain.RecipeFilter, page int) ([]*domain.Recipe, error) {
	//allIngredients := 0
	//if filter.AvailableIngredients == nil || len(filter.AvailableIngredients) == 0 {
	//	allIngredients = 1
	//}
	//allTypes := 0
	//if filter.SaladTypes == nil || len(filter.SaladTypes) == 0 {
	//	allTypes = 1
	//}
	//
	//type ingredientRes struct {
	//	recipeId uuid.UUID `gorm:"column:recipeId"`
	//	matches  int       `gorm:"column:matches"`
	//}
	//var ingredientsMatches ingredientRes
	//
	//err := r.db.WithContext(ctx).
	//	Model(&rDomain.IngredientLink{}).
	//	Select("recipeId, count(*) as matches").
	//	Where("ingredientId = ?", filter.AvailableIngredients).
	//	Or(r.db.Where("? = 1", allIngredients)).
	//	Group("recipeId").
	//	First(&ingredientsMatches).Error
	//if err != nil {
	//	return nil, fmt.Errorf("getting all recipes: %w", err)
	//}

	ingredientUUIDS := ""
	if filter.AvailableIngredients != nil && len(filter.AvailableIngredients) != 0 {
		ingredientUUIDS = uuidsToString(filter.AvailableIngredients)
	} else {
		ingredientUUIDS = "ingredientId"
	}

	saladTypesUUIDS := ""
	if filter.SaladTypes != nil && len(filter.SaladTypes) != 0 {
		saladTypesUUIDS = uuidsToString(filter.SaladTypes)
	} else {
		saladTypesUUIDS = "typeId"
	}

	query := `with matchCounts as (select recipeId, count(*) as matches
		from saladRecipes.recipeIngredient
		where ingredientId in (?)
		group by recipeId),
	totalIngredients as (select recipeId, count(*) as ingredientsCount
		from saladRecipes.recipeIngredient
		group by recipeId),
	availableRecipes as (select id, saladId, status, numberOfServings, timeToCook, rating
		from saladRecipes.recipe
		where id in (select matchCounts.recipeId
             from matchCounts join totalIngredients on matchCounts.recipeId = totalIngredients.recipeId
             where matches = ingredientsCount)),
	requestedTypes as (select saladId
		from saladRecipes.typesOfSalads
		where typeId in (?))
		group by saladId)
	select id, saladId, status, numberOfServings, timeToCook, rating
		from availableRecipes
		where
			saladId in (select saladId from requestedTypes) and
			(rating is null or rating > ?)
		order by case
    		when rating is null then rating
    		else 0 end, rating desc 
		offset ?
		limit ?`

	var dbRecipes []*rDomain.Recipe
	err := r.db.WithContext(ctx).
		Raw(query, ingredientUUIDS, saladTypesUUIDS, filter.MinRate, PageSize*(page-1), PageSize*(page+1)).
		Scan(&dbRecipes).Error
	if err != nil {
		return nil, fmt.Errorf("getting recipes: %w", err)
	}
	recipes := make([]*domain.Recipe, 0)
	for _, recipe := range dbRecipes {
		recipes = append(recipes, rDomain.ToRecipeBL(recipe))
	}
	return recipes, nil
}

func (r *recipeRepository) Update(ctx context.Context, recipe *domain.Recipe) error {
	dbRecipe := rDomain.ToRecipeDB(recipe)
	err := r.db.WithContext(ctx).
		Save(&dbRecipe).Error
	if err != nil {
		return fmt.Errorf("updating recipe: %w", err)
	}
	return nil
}

func (r *recipeRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&rDomain.Recipe{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting recipe by id: %w", err)
	}
	return nil
}
