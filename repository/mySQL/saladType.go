package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type saladTypeRepository struct {
	db *gorm.DB
}

func NewSaladTypeRepository(db *gorm.DB) rDomain.ISaladTypeRepository {
	return &saladTypeRepository{
		db: db,
	}
}

func (r *saladTypeRepository) Create(ctx context.Context, saladType *domain.SaladType) error {
	typeDb := rDomain.ToSaladTypeDB(saladType)
	typeDb.ID = uuid.New()
	err := r.db.WithContext(ctx).
		Create(&typeDb).Error
	if err != nil {
		return fmt.Errorf("creating salad type: %w", err)
	}
	return nil
}

func (r *saladTypeRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.SaladType, error) {
	var typeDb rDomain.SaladType
	err := r.db.WithContext(ctx).
		First(&typeDb, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting salad type by id: %w", err)
	}
	return rDomain.ToSaladTypeBL(&typeDb), nil
}

func (r *saladTypeRepository) GetAll(ctx context.Context, page int) ([]*domain.SaladType, int, error) {
	var saladTypes []*rDomain.SaladType
	err := r.db.WithContext(ctx).
		Limit(PageSize).
		Offset(PageSize * (page - 1)).
		Find(&saladTypes).Error
	if err != nil {
		return nil, 0, fmt.Errorf("getting salad types: %w", err)
	}

	resSaladTypes := make([]*domain.SaladType, 0)
	for _, saladType := range saladTypes {
		resSaladTypes = append(resSaladTypes, rDomain.ToSaladTypeBL(saladType))
	}

	var tmp []*rDomain.SaladType
	count := r.db.WithContext(ctx).
		Find(&tmp).RowsAffected
	pagesCount := count / PageSize
	if count%PageSize != 0 {
		pagesCount++
	}

	return resSaladTypes, int(pagesCount), nil
}

func (r *saladTypeRepository) GetAllBySaladId(ctx context.Context, saladId uuid.UUID) ([]*domain.SaladType, error) {
	var saladTypes []*rDomain.SaladType
	var saladIds []uuid.UUID

	err := r.db.WithContext(ctx).
		Table("typesOfSalads").
		Where("saladId = ?", saladId).
		Select("typeId").Scan(&saladIds).Error
	if err != nil {
		return nil, fmt.Errorf("getting salad type ids: %w", err)
	}

	if len(saladIds) != 0 {
		err = r.db.WithContext(ctx).
			Find(&saladTypes, saladIds).Error
	}
	if err != nil {
		return nil, fmt.Errorf("getting salad types: %w", err)
	}

	resSaladTypes := make([]*domain.SaladType, 0)
	for _, saladType := range saladTypes {
		resSaladTypes = append(resSaladTypes, rDomain.ToSaladTypeBL(saladType))
	}

	return resSaladTypes, nil
}

func (r *saladTypeRepository) Update(ctx context.Context, saladType *domain.SaladType) error {
	dbSaladType := rDomain.ToSaladTypeDB(saladType)
	err := r.db.WithContext(ctx).
		Save(dbSaladType).Error
	if err != nil {
		return fmt.Errorf("updating salad type: %w", err)
	}
	return nil
}

func (r *saladTypeRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&rDomain.SaladType{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting salad type by id: %w", err)
	}
	return nil
}

func (r *saladTypeRepository) Link(ctx context.Context, saladId uuid.UUID, saladTypeId uuid.UUID) error {
	dbLink := rDomain.TypeLink{
		ID:      uuid.New(), // TODO: mb should generate uuid there
		SaladId: saladId,
		TypeId:  saladTypeId,
	}
	err := r.db.WithContext(ctx).
		Create(&dbLink).Error
	if err != nil {
		return fmt.Errorf("linking types from salad: %w", err)
	}
	return nil
}

func (r *saladTypeRepository) Unlink(ctx context.Context, saladId uuid.UUID, saladTypeId uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("saladId = ?", saladId).
		Where("typeId = ?", saladTypeId).
		Delete(&rDomain.TypeLink{}).Error
	if err != nil {
		return fmt.Errorf("unlinking types from salad: %w", err)
	}
	return nil
}
