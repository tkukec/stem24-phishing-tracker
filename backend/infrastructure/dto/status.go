package dto

import (
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"time"
)

// NewStatus constructor for Status
func NewStatus(status *models.Status) *Status {
	return &Status{
		ID:        status.ID,
		CreatedAt: status.CreatedAt,
		UpdatedAt: status.UpdatedAt,
		DeletedAt: status.DeletedAt.Time,
		Name:      status.Name,
	}
}

// Status ....
type Status struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Name      string    `json:"names"`
}
