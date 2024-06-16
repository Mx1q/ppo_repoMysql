package domain

import (
	"context"
	"github.com/Mx1q/ppo_services/domain"
	"github.com/google/uuid"
)

type Comment struct {
	ID       uuid.UUID `gorm:"primaryKey"`
	AuthorID uuid.UUID `gorm:"column:author"`
	SaladID  uuid.UUID `gorm:"column:salad"`
	Text     string    `gorm:"column:text"`
	Rating   int       `gorm:"column:rating"`
}

func (Comment) TableName() string {
	return "comment"
}

func ToCommentDB(comment *domain.Comment) *Comment {
	return &Comment{
		ID:       comment.ID,
		AuthorID: comment.AuthorID,
		SaladID:  comment.SaladID,
		Text:     comment.Text,
		Rating:   comment.Rating,
	}
}

func ToCommentBL(comment *Comment) *domain.Comment {
	return &domain.Comment{
		ID:       comment.ID,
		AuthorID: comment.AuthorID,
		SaladID:  comment.SaladID,
		Text:     comment.Text,
		Rating:   comment.Rating,
	}
}

type ICommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	GetById(ctx context.Context, id uuid.UUID) (*domain.Comment, error)
	GetBySaladAndUser(ctx context.Context, saladId uuid.UUID, userId uuid.UUID) (*domain.Comment, error)
	GetAllBySaladID(ctx context.Context, saladId uuid.UUID, page int) ([]*domain.Comment, int, error)
	Update(ctx context.Context, comment *domain.Comment) error
	DeleteById(ctx context.Context, id uuid.UUID) error
}
