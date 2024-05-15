package repositories

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
)

// NewComment constructor for Comment
func NewComment(conn *database.Connection) CommentRepository {
	return &Comment{
		conn: conn,
	}
}

type CommentRepository interface {
	Persist(tenantID string, comment *models.Comment) (*models.Comment, error)
	Update(tenantID string, comment *models.Comment) (*models.Comment, error)
	Delete(tenantID string, comment *models.Comment) error
	GetByEmail(tenantID, email string) (*models.Comment, error)
	Get(tenantID string, ID string) (*models.Comment, error)
	GetAll(tenantID string, query database.Query, with ...string) ([]*models.Comment, error)
}

// Comment ....
type Comment struct {
	conn *database.Connection
}

// Persist ....
func (r *Comment) Persist(tenantID string, comment *models.Comment) (*models.Comment, error) {
	comment.TenantID = tenantID
	if err := r.conn.GetConnectionWithPreload([]string{"SkillGroups"}).Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// Update ....
func (r *Comment) Update(tenantID string, comment *models.Comment) (*models.Comment, error) {
	comment.TenantID = tenantID
	if err := r.conn.GetConnectionWithPreload([]string{"SkillGroups"}).Save(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// Delete ....
func (r *Comment) Delete(tenantID string, comment *models.Comment) error {
	comment.TenantID = tenantID
	if err := r.conn.GetConnection().Delete(comment).Error; err != nil {
		return err
	}
	return nil
}

func (r *Comment) GetByEmail(tenantID, email string) (*models.Comment, error) {
	var comment models.Comment
	if err := r.conn.GetConnectionWithPreload(nil).Where("tenant_id = ?", tenantID).First(&comment, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// Get ....
func (r *Comment) Get(tenantID string, ID string) (*models.Comment, error) {
	var comment models.Comment
	if err := r.conn.GetConnectionWithPreload([]string{"SkillGroups"}).Where("tenant_id = ?", tenantID).First(&comment, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetAll ...
func (r *Comment) GetAll(tenantID string, query database.Query, with ...string) ([]*models.Comment, error) {
	var records []*models.Comment
	preload := []string{"SkillGroups"}
	if len(with) > 0 {
		preload = append(preload, with...)
	}

	tx := r.conn.GetConnectionWithPreload(preload)

	if query != nil {
		if query.Limit() != 0 {
			tx.Limit(query.Limit())
		}

		if query.Offset() != 0 {
			tx.Offset(query.Offset())
		}

		tx.Order(query.OrderBy())

		for _, item := range query.Build() {
			tx.Where(fmt.Sprintf("%s %s ?", item.Key(), item.Operator()), item.Value())
		}
	}

	tx.Model(&models.Comment{}).
		Distinct().
		Where("tenant_id = ?", tenantID)

	if err := tx.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
