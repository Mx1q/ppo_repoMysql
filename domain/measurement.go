package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

type Measurement struct {
	ID    uuid.UUID `gorm:"primaryKey"`
	Name  string    `gorm:"column:name"`
	Grams int       `gorm:"column:grams"`
}

type MeasurementLink struct {
	Measurement uuid.UUID
	Amount      int
}

func (Measurement) TableName() string {
	return "measurement"
}

func ToMeasurementDB(measurement *domain.Measurement) *Measurement {
	return &Measurement{
		ID:    measurement.ID,
		Name:  measurement.Name,
		Grams: measurement.Grams,
	}
}

func ToMeasurementBL(measurement *Measurement) *domain.Measurement {
	return &domain.Measurement{
		ID:    measurement.ID,
		Name:  measurement.Name,
		Grams: measurement.Grams,
	}
}

type IMeasurementRepository interface {
	Create(ctx context.Context, measurement *domain.Measurement) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.Measurement, error)
	GetByRecipeId(ctx context.Context, ingredientId uuid.UUID, recipeId uuid.UUID) (*domain.Measurement, int, error)
	GetAll(ctx context.Context) ([]*domain.Measurement, error)
	UpdateLink(ctx context.Context, linkId uuid.UUID, measurementId uuid.UUID, amount int) error
	Update(ctx context.Context, measurement *domain.Measurement) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
