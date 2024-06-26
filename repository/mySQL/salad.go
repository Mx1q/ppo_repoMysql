package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"slices"
)

type saladRepository struct {
	db *gorm.DB
}

func NewSaladRepository(db *gorm.DB) rDomain.ISaladRepository {
	return &saladRepository{
		db: db,
	}
}

func (r *saladRepository) Create(ctx context.Context, salad *domain.Salad) (uuid.UUID, error) {
	dbSalad := rDomain.ToSaladDB(salad)
	dbSalad.ID = uuid.New()
	err := r.db.WithContext(ctx).
		Create(&dbSalad).Error

	if err != nil {
		return uuid.Nil, fmt.Errorf("creating salad: %w", err)
	}
	return dbSalad.ID, nil
}

func (r *saladRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Salad, error) {
	var salad rDomain.Salad
	err := r.db.WithContext(ctx).
		First(&salad, id).Error

	if err != nil {
		return nil, fmt.Errorf("getting salad by id: %w", err)
	}
	return rDomain.ToSaladBL(&salad), nil
}

func (r *saladRepository) GetAll(ctx context.Context, filter *domain.RecipeFilter, page int) ([]*domain.Salad, int, error) {
	type iRes struct {
		id    uuid.UUID `gorm:"column:id"`
		count int       `gorm:"column:cnt"`
	}

	var availableIngredients []*iRes
	var totalCount []*iRes
	var filteredIngredients []uuid.UUID

	tiRows, err := r.db.WithContext(ctx).
		Table("recipe").
		Select("recipe.saladId as id", "count(*) as cnt").
		Joins("left join recipeIngredient on recipe.id = recipeIngredient.recipeId").
		Group("recipe.id").
		Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("fetching salads: %w", err)
	}
	defer tiRows.Close()
	for tiRows.Next() {
		var id uuid.UUID
		var count int
		tiRows.Scan(&id, &count)
		totalCount = append(totalCount, &iRes{id: id, count: count})
	}

	if filter.AvailableIngredients == nil || len(filter.AvailableIngredients) == 0 {
		for _, cnt := range totalCount {
			filteredIngredients = append(filteredIngredients, cnt.id)
		}
	} else {
		iRows, err := r.db.WithContext(ctx).
			Table("recipe").
			Select("recipe.saladId as id", "count(*) as cnt").
			Where("recipeIngredient.ingredientId in ?", filter.AvailableIngredients).
			Joins("left join recipeIngredient on recipe.id = recipeIngredient.recipeId").
			Group("recipe.id").
			Rows()
		if err != nil {
			return nil, 0, fmt.Errorf("fetching salads: %w", err)
		}
		defer iRows.Close()
		for iRows.Next() {
			var id uuid.UUID
			var count int
			iRows.Scan(&id, &count)
			availableIngredients = append(availableIngredients, &iRes{id: id, count: count})
		}

		for i := 0; i < len(totalCount); i++ {
			for _, cnt := range availableIngredients {
				if cnt.id == totalCount[i].id {
					if cnt.count == totalCount[i].count {
						filteredIngredients = append(filteredIngredients, totalCount[i].id)
					}
					break
				}
			}
		}
	}

	var tIds []uuid.UUID
	err = r.db.WithContext(ctx).
		Table("recipe").
		Select("recipe.saladId as id").
		Where("typesOfSalads.typeId in ?", filter.SaladTypes).
		Joins("left join typesOfSalads on recipe.saladId = typesOfSalads.saladId").
		Group("recipe.id").
		Scan(&tIds).Error
	if err != nil {
		return nil, 0, fmt.Errorf("fetching salads: %w", err)
	}

	var twiceSorted []uuid.UUID
	if filter.SaladTypes == nil || len(filter.SaladTypes) == 0 {
		for _, f := range filteredIngredients {
			twiceSorted = append(twiceSorted, f)
		}
	} else {
		for _, id := range tIds {
			if slices.Contains(filteredIngredients, id) {
				twiceSorted = append(twiceSorted, id)
			}
		}
	}

	var dbSalads []*rDomain.Salad
	err = r.db.WithContext(context.Background()).
		Table("salad").
		Select("salad.id", "salad.authorId", "salad.name", "salad.description").
		Joins("left join recipe on recipe.saladId = salad.id").
		Where("salad.id in ?", twiceSorted).
		Where(
			r.db.WithContext(ctx).
				Where("recipe.rating >= ?", filter.MinRate).Or("recipe.rating is null"),
		).
		Where("recipe.status = ?", filter.Status).
		Order("recipe.rating desc").
		Limit(PageSize).
		Offset(PageSize * (page - 1)).
		Find(&dbSalads).Error

	salads := make([]*domain.Salad, 0)
	for _, salad := range dbSalads {
		salads = append(salads, rDomain.ToSaladBL(salad))
	}

	var tmp []*rDomain.Salad
	count := r.db.WithContext(ctx).
		Find(&tmp).RowsAffected
	numPages := count / PageSize
	if count%PageSize != 0 {
		numPages++
	}

	return salads, int(numPages), nil
}

func (r *saladRepository) GetAllByUserId(ctx context.Context, id uuid.UUID) ([]*domain.Salad, error) {
	var dbSalads []*rDomain.Salad
	err := r.db.WithContext(ctx).
		Table("salad").
		Where("authorId = ?", id).
		Scan(&dbSalads).Error
	if err != nil {
		return nil, fmt.Errorf("getting salads by user id: %w", err)
	}

	salads := make([]*domain.Salad, 0)
	for _, salad := range dbSalads {
		salads = append(salads, rDomain.ToSaladBL(salad))
	}
	return salads, nil
}

func (r *saladRepository) GetAllRatedByUser(ctx context.Context, userId uuid.UUID, page int) ([]*domain.Salad, int, error) {
	var saladIds []uuid.UUID
	rows, err := r.db.WithContext(ctx).
		Table("comment").
		Select("comment.salad").
		Where("author = ?", userId).
		Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("fetching salads: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		rows.Scan(&id)
		saladIds = append(saladIds, id)
	}

	var dbSalads []*rDomain.Salad
	if len(saladIds) != 0 {
		err := r.db.WithContext(ctx).
			Table("salad").
			Limit(PageSize).
			Offset(PageSize*(page-1)).
			Find(&dbSalads, saladIds).Error

		if err != nil {
			return nil, 0, fmt.Errorf("getting salads rated by user: %w", err)
		}
	}

	salads := make([]*domain.Salad, 0)
	for _, salad := range dbSalads {
		salads = append(salads, rDomain.ToSaladBL(salad))
	}

	numPages := 0
	if len(saladIds) != 0 {
		var tmp []*rDomain.Salad
		count := r.db.WithContext(ctx).
			Find(&tmp, saladIds).RowsAffected
		numPages = int(count / PageSize)
		if count%PageSize != 0 {
			numPages++
		}
	}

	return salads, int(numPages), nil
}

func (r *saladRepository) Update(ctx context.Context, salad *domain.Salad) error {
	dbSalad := rDomain.ToSaladDB(salad)
	err := r.db.WithContext(ctx).
		Save(&dbSalad).Error
	if err != nil {
		return fmt.Errorf("updating salad: %w", err)
	}
	return nil
}

func (r *saladRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&rDomain.Salad{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting salad by id: %w", err)
	}
	return nil
}
