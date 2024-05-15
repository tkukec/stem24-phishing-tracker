package dto

import (
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"time"
)

func NewComments(comments []models.Comment) []*Comment {
	dtoComments := make([]*Comment, 0)
	for _, v := range comments {
		dtoComments = append(dtoComments, NewComment(&v))
	}
	return dtoComments
}

// NewComment constructor for Comment
func NewComment(comment *models.Comment) *Comment {
	return &Comment{
		ID:          comment.ID,
		CreatedAt:   comment.CreatedAt,
		UpdatedAt:   comment.UpdatedAt,
		DeletedAt:   comment.DeletedAt.Time,
		Description: comment.Description,
	}
}

// Comment ....
type Comment struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	Description string    `json:"description"`
}
