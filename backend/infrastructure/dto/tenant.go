package dto

import (
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"time"
)

// NewTenant constructor for Tenant
func NewTenant(tenant *models.Tenant) *Tenant {
	return &Tenant{
		ID:        tenant.ID,
		CreatedAt: tenant.CreatedAt,
		UpdatedAt: tenant.UpdatedAt,
		DeletedAt: tenant.DeletedAt.Time,
		Name:      tenant.Name,
	}
}

// Tenant ....
type Tenant struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Name      string    `json:"names"`
}
