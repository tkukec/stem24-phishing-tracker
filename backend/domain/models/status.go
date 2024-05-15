package models

const (
	StatusModelName = "status"
)

type Status struct {
	Model
	Name string `gorm:"type:varchar(255);not null"`
}

func NewStatus(name string) *Status {
	return &Status{
		Name: name,
	}
}
