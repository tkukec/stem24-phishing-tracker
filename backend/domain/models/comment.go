package models

const (
	CommentModelName = "comment"
)

type Comment struct {
	Model
	Description string `gorm:"type:varchar(1500);not null"`
	EventID     string
	Event       Event
}

func NewComment(description string) *Comment {
	return &Comment{
		Description: description,
	}
}
