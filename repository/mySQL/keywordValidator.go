package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type keywordValidatorRepository struct {
	db *gorm.DB
}

func NewKeywordValidatorRepository(db *gorm.DB) rDomain.IKeywordValidatorRepository {
	return &keywordValidatorRepository{
		db: db,
	}
}

func (r *keywordValidatorRepository) Create(ctx context.Context, word *domain.KeyWord) error {
	dbKeyWord := rDomain.ToKeyWordDB(word)
	dbKeyWord.ID = uuid.New()
	err := r.db.WithContext(ctx).
		Create(&dbKeyWord).Error
	if err != nil {
		return fmt.Errorf("creating keyword: %w", err)
	}
	return nil
}

func (r *keywordValidatorRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.KeyWord, error) {
	var keyWord rDomain.KeyWord
	err := r.db.WithContext(ctx).
		First(&keyWord, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting keyword by id: %w", err)
	}
	return rDomain.ToKeyWordBL(&keyWord), nil
}

func (r *keywordValidatorRepository) GetAll(ctx context.Context) (map[string]uuid.UUID, error) {
	var keyWords []*rDomain.KeyWord
	err := r.db.WithContext(ctx).
		Find(&keyWords).Error
	if err != nil {
		return nil, fmt.Errorf("getting all keywords: %w", err)
	}

	resKeywords := make(map[string]uuid.UUID)
	for _, word := range keyWords {
		resKeywords[word.Word] = word.ID
	}
	return resKeywords, nil
}

func (r *keywordValidatorRepository) Update(ctx context.Context, word *domain.KeyWord) error {
	dbKeyword := rDomain.ToKeyWordDB(word)
	err := r.db.WithContext(ctx).
		Save(&dbKeyword).Error
	if err != nil {
		return fmt.Errorf("updating keyword: %w", err)
	}
	return nil
}

func (r *keywordValidatorRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&rDomain.KeyWord{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting keyword by id: %w", err)
	}
	return nil
}
