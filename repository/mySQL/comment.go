package mysql

import (
	"context"
	"fmt"
	rDomain "github.com/Mx1q/ppo_repoMysql/domain"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) rDomain.ICommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	dbComment := rDomain.ToCommentDB(comment)
	dbComment.ID = uuid.New()
	err := r.db.WithContext(ctx).
		Create(&dbComment).Error
	if err != nil {
		return fmt.Errorf("creating comment: %w", err)
	}
	return nil
}

func (r *commentRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Comment, error) {
	var dbComment rDomain.Comment
	err := r.db.WithContext(ctx).
		First(&dbComment, id).Error
	if err != nil {
		return nil, fmt.Errorf("getting comment by id: %w", err)
	}
	return rDomain.ToCommentBL(&dbComment), nil
}

func (r *commentRepository) GetBySaladAndUser(ctx context.Context, saladId uuid.UUID, userId uuid.UUID) (*domain.Comment, error) {
	var dbComment rDomain.Comment
	err := r.db.WithContext(ctx).
		Model(&rDomain.Comment{}).
		Where("salad = ?", saladId).
		Where("author = ?", userId).
		First(&dbComment).Error
	if err != nil {
		return nil, fmt.Errorf("getting comment by salad and user IDs: %w", err)
	}
	return rDomain.ToCommentBL(&dbComment), nil
}

func (r *commentRepository) GetAllBySaladID(ctx context.Context, saladId uuid.UUID, page int) ([]*domain.Comment, int, error) {
	var comments []*rDomain.Comment
	err := r.db.WithContext(ctx).
		Where("salad = ?", saladId).
		Limit(PageSize).
		Offset(PageSize * (page - 1)).
		Find(&comments).Error
	if err != nil {
		return nil, 0, fmt.Errorf("getting comments by salad id: %w", err)
	}

	resComments := make([]*domain.Comment, 0)
	for _, comment := range comments {
		resComments = append(resComments, rDomain.ToCommentBL(comment))
	}

	var tmp []*rDomain.Comment
	count := r.db.WithContext(ctx).
		Find(&tmp).RowsAffected
	numPages := count / PageSize
	if count%PageSize != 0 {
		numPages++
	}

	return resComments, int(numPages), nil
}

func (r *commentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	dbComment := rDomain.ToCommentDB(comment)
	err := r.db.WithContext(ctx).
		Save(&dbComment).Error
	if err != nil {
		return fmt.Errorf("updating comment: %w", err)
	}
	return nil
}

func (r *commentRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Delete(&rDomain.Comment{}, id).Error
	if err != nil {
		return fmt.Errorf("deleting comment by id: %w", err)
	}
	return nil
}
