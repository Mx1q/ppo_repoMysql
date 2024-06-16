package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type recipeStepRepository struct {
	db *gorm.DB
}

func NewRecipeStepRepository(db *gorm.DB) rDomain.IRecipeStepRepository {
	return &recipeStepRepository{
		db: db,
	}
}

func (r *recipeStepRepository) Create(ctx context.Context, recipeStep *domain.RecipeStep) error {
	query := `insert into saladRecipes.recipeStep(name, description, recipeId, stepNum)
	values (@name, @description, @recipe, 
        (select freeStep from (select case
			when max(stepNum) is null then 1
			else max(stepNum) + 1
			end as freeStep
		from saladRecipes.recipeStep
		where recipeId = @recipe) as tmp))`

	err := r.db.WithContext(ctx).
		Exec(query, map[string]interface{}{
			"recipe":      recipeStep.RecipeID,
			"name":        recipeStep.Name,
			"description": recipeStep.Description,
		}).Error
	if err != nil {
		return fmt.Errorf("creating recipe step: %w", err)
	}
	return nil
}

func (r *recipeStepRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.RecipeStep, error) {
	var dbStep rDomain.RecipeStep
	err := r.db.WithContext(ctx).
		First(&dbStep, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting recipe step by id: %w", err)
	}
	return rDomain.ToStepBL(&dbStep), nil
}

func (r *recipeStepRepository) GetAllByRecipeID(ctx context.Context, recipeId uuid.UUID) ([]*domain.RecipeStep, error) {
	var dbSteps []*rDomain.RecipeStep
	err := r.db.WithContext(ctx).
		Table("recipeStep").
		Order("stepNum").
		Where("recipeId = ?", recipeId).
		Scan(&dbSteps).Error
	if err != nil {
		return nil, fmt.Errorf("getting all recipe steps: %w", err)
	}

	steps := make([]*domain.RecipeStep, 0)
	for _, step := range dbSteps {
		steps = append(steps, rDomain.ToStepBL(step))
	}
	return steps, nil
}

func (r *recipeStepRepository) Update(ctx context.Context, recipeStep *domain.RecipeStep) error {
	tx := r.db.WithContext(ctx).
		Begin()
	tx.SavePoint("rollbackPoint")
	err := tx.Error
	defer func() {
		if err != nil {
			rollbackErr := tx.RollbackTo("rollbackPoint").Error
			if rollbackErr != nil {
				err = fmt.Errorf("%v: %w", rollbackErr, err)
			}
		}
	}()

	type maxRes struct {
		maxNum  int
		stepNum int
	}
	var res maxRes

	err = tx.WithContext(context.Background()).
		Table("recipeStep").
		Where("recipeId = ?", recipeStep.RecipeID).
		Select("case when max(stepNum) is null then 0 else max(stepNum) end as maxNum").
		Scan(&res.maxNum).Error
	if err != nil {
		return fmt.Errorf("updating recipe step (checking max step num): %w", err)
	}

	err = tx.WithContext(context.Background()).
		Table("recipeStep").
		Where("id = ?", recipeStep.ID).
		Select("stepNum").
		Scan(&res.stepNum).Error
	if err != nil {
		return fmt.Errorf("updating recipe step (checking max step num): %w", err)
	}
	if recipeStep.StepNum > res.maxNum {
		return fmt.Errorf("updating recipe step: step num out of range")
	}

	if recipeStep.StepNum < res.stepNum {
		err = tx.WithContext(ctx).
			Table("recipeStep").
			Where("recipeId = ?", recipeStep.RecipeID).
			Where("stepNum between ? and ?", recipeStep.StepNum, res.stepNum-1).
			Update("stepNum", gorm.Expr("stepNum + 1")).Error
	} else {
		err = tx.WithContext(ctx).
			Table("recipeStep").
			Where("recipeId = ?", recipeStep.RecipeID).
			Where("stepNum between ? and ?", res.stepNum+1, recipeStep.StepNum).
			Update("stepNum", gorm.Expr("stepNum - 1")).Error
	}
	if err != nil {
		return fmt.Errorf("updating recipe step (moving other steps): %w", err)
	}

	query := `update saladRecipes.recipeStep
		set
			name = @name,
			description = @description,
			stepNum = @stepNum
		where id = @id`
	err = tx.WithContext(ctx).
		Exec(query, map[string]interface{}{
			"name":        recipeStep.Name,
			"description": recipeStep.Description,
			"stepNum":     recipeStep.StepNum,
			"id":          recipeStep.ID,
		}).Error
	if err != nil {
		return fmt.Errorf("updating recipe step: %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("updating recipe step (commiting transaction): %w", err)
	}
	return nil
}

func (r *recipeStepRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	tx := r.db.WithContext(ctx).
		Begin()
	tx.SavePoint("rollbackPoint")
	err := tx.Error
	defer func() {
		if err != nil {
			rollbackErr := tx.RollbackTo("rollbackPoint").Error
			if rollbackErr != nil {
				err = fmt.Errorf("%v: %w", rollbackErr, err)
			}
		}
	}()

	type queryRes struct {
		recipeId uuid.UUID
		stepNum  int
	}
	var res queryRes

	var dbStep rDomain.RecipeStep
	err = tx.WithContext(ctx).
		Find(&dbStep, id).Error
	res.stepNum = dbStep.StepNum
	res.recipeId = dbStep.RecipeID

	if err != nil {
		return fmt.Errorf("deleting recipe step by id (getting recipe ID): %w", err)
	}

	fmt.Println("!!!!!RES!!!!!", res.recipeId, res.stepNum)

	err = tx.WithContext(ctx).
		Table("recipeStep").
		Where("recipeId = ?", res.recipeId).
		Where("stepNum > ?", res.stepNum).
		Update("stepNum", gorm.Expr("stepNum - 1")).Error

	if err != nil {
		return fmt.Errorf("updating recipe step (moving other steps): %w", err)
	}
	err = tx.WithContext(ctx).
		Delete(&rDomain.RecipeStep{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting recipe step by id: %w", err)
	}

	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("deleting recipe step (commiting transaction): %w", err)
	}
	return nil
}

func (r *recipeStepRepository) DeleteAllByRecipeID(ctx context.Context, recipeId uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("recipeId = ?", recipeId).
		Delete(&rDomain.RecipeStep{}).Error
	if err != nil {
		return fmt.Errorf("deleting all recipe steps by id: %w", err)
	}
	return nil
}
