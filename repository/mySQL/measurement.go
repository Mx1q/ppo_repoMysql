package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const PageSize = 30

type measurementRepository struct {
	db *gorm.DB
}

func NewMeasrementRepository(db *gorm.DB) rDomain.IMeasurementRepository {
	return &measurementRepository{
		db: db,
	}
}

func (r *measurementRepository) Create(ctx context.Context, measurement *domain.Measurement) error {
	dbMeasurement := rDomain.ToMeasurementDB(measurement)
	dbMeasurement.ID = uuid.New()
	err := r.db.WithContext(ctx).
		Create(&dbMeasurement).Error
	if err != nil {
		return fmt.Errorf("creating measurement: %w", err)
	}
	return nil
}

func (r *measurementRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Measurement, error) {
	var measurement rDomain.Measurement
	err := r.db.WithContext(ctx).
		First(&measurement, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting measurement by id: %w", err)
	}
	return rDomain.ToMeasurementBL(&measurement), nil
}

func (r *measurementRepository) GetByRecipeId(ctx context.Context,
	ingredientId uuid.UUID, recipeId uuid.UUID) (*domain.Measurement, int, error) {
	var link rDomain.MeasurementLink
	err := r.db.WithContext(ctx).
		Table("recipeIngredient").
		Where("ingredientId = ?", ingredientId).
		Where("recipeId = ?", recipeId).
		Select("measurement, amount").
		Scan(&link).Error
	if err != nil {
		return nil, 0, fmt.Errorf("getting measurement link: %w", err)
	}

	var measurement rDomain.Measurement
	err = r.db.WithContext(ctx).
		First(&measurement, link.Measurement).Error
	if err != nil {
		return nil, 0, fmt.Errorf("getting measurement by recipe and ingredient: %w", err)
	}
	return rDomain.ToMeasurementBL(&measurement), link.Amount, nil
}

func (r *measurementRepository) GetAll(ctx context.Context) ([]*domain.Measurement, error) {
	var measurements []*rDomain.Measurement
	err := r.db.WithContext(ctx).
		Find(&measurements).Error
	if err != nil {
		return nil, fmt.Errorf("getting measurements: %w", err)
	}

	resMeasurements := make([]*domain.Measurement, 0)
	for _, measurement := range measurements {
		resMeasurements = append(resMeasurements, rDomain.ToMeasurementBL(measurement))
	}
	return resMeasurements, nil
}

func (r *measurementRepository) Update(ctx context.Context, measurement *domain.Measurement) error {
	dbMeasurement := rDomain.ToMeasurementDB(measurement)
	err := r.db.WithContext(ctx).
		Save(&dbMeasurement).Error
	if err != nil {
		return fmt.Errorf("updating measurement: %w", err)
	}
	return nil
}

func (r *measurementRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&rDomain.Measurement{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting measurement by id: %w", err)
	}
	return nil
}

func (r *measurementRepository) UpdateLink(ctx context.Context, linkId uuid.UUID, measurementId uuid.UUID, amount int) error {
	err := r.db.WithContext(ctx).
		Raw("update saladRecipes.recipeIngredient set measurement = ?, amount = ? where id = ?", measurementId, amount, linkId).
		Error
	if err != nil {
		return fmt.Errorf("updating measurement: %w", err)
	}
	return nil
}
