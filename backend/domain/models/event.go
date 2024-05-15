package models

import "time"

const (
	EventModelName = "event"
)

type Event struct {
	Model
	Name             string `gorm:"type:varchar(500);not null"`
	Date             time.Time
	Brand            string `gorm:"type:varchar(500);not null"`
	Description      string `gorm:"type:varchar(1500)"`
	MalURL           string `gorm:"type:varchar(1500)"`
	MalDomainRegDate time.Time
	DNSRecord        string   `gorm:"type:varchar(500)"`
	Keywords         []string `gorm:"type:varchar(1500)"`
	Status           Status
	Comments         []Comment
}
