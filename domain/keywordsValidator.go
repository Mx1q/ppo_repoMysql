package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

type KeyWord struct {
	ID   uuid.UUID `gorm:"primaryKey"`
	Word string    `gorm:"column:word"`
}

func (KeyWord) TableName() string {
	return "word"
}

func ToKeyWordDB(keyWord *domain.KeyWord) *KeyWord {
	return &KeyWord{
		ID:   keyWord.ID,
		Word: keyWord.Word,
	}
}

func ToKeyWordBL(keyWord *KeyWord) *domain.KeyWord {
	return &domain.KeyWord{
		ID:   keyWord.ID,
		Word: keyWord.Word,
	}
}

type IKeywordValidatorRepository interface {
	Create(ctx context.Context, word *domain.KeyWord) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.KeyWord, error)
	GetAll(ctx context.Context) (map[string]uuid.UUID, error)
	Update(ctx context.Context, word *domain.KeyWord) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
