package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const (
	TenantModelName = "tenant"
)

type Tenant struct {
	ID        string `gorm:"type:char(36);primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name string `gorm:"unique"`
}

func TenantSeed() []*Tenant {
	return []*Tenant{
		&Tenant{
			Name: "*",
		},
	}
}

func (m *Tenant) BeforeCreate(scope *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	return nil
}
