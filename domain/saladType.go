package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

type SaladType struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
}

type TypeLink struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	SaladId uuid.UUID `gorm:"column:saladId"`
	TypeId  uuid.UUID `gorm:"column:typeId"`
}

func (SaladType) TableName() string {
	return "saladType"
}

func (TypeLink) TableName() string {
	return "typesOfSalads"
}

func ToSaladTypeDB(saladType *domain.SaladType) *SaladType {
	return &SaladType{
		ID:          saladType.ID,
		Name:        saladType.Name,
		Description: saladType.Description,
	}
}

func ToSaladTypeBL(saladType *SaladType) *domain.SaladType {
	return &domain.SaladType{
		ID:          saladType.ID,
		Name:        saladType.Name,
		Description: saladType.Description,
	}
}

type ISaladTypeRepository interface {
	Create(ctx context.Context, saladType *domain.SaladType) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.SaladType, error)
	GetAll(ctx context.Context, page int) ([]*domain.SaladType, int, error)
	GetAllBySaladId(ctx context.Context, saladId uuid.UUID) ([]*domain.SaladType, error)
	Update(ctx context.Context, saladType *domain.SaladType) error
	Link(ctx context.Context, saladId uuid.UUID, saladTypeId uuid.UUID) error
	Unlink(ctx context.Context, saladId uuid.UUID, saladTypeId uuid.UUID) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
