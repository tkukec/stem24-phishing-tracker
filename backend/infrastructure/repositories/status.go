package repositories

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
)

// NewStatus constructor for Status
func NewStatus(conn *database.Connection) StatusRepository {
	return &Status{
		conn: conn,
	}
}

type StatusRepository interface {
	Persist(tenantID string, status *models.Status) (*models.Status, error)
	Update(tenantID string, status *models.Status) (*models.Status, error)
	Delete(tenantID string, status *models.Status) error
	GetByEmail(tenantID, email string) (*models.Status, error)
	Get(tenantID string, ID string) (*models.Status, error)
	GetAll(tenantID string, query database.Query, with ...string) ([]*models.Status, error)
}

// Status ....
type Status struct {
	conn *database.Connection
}

// Persist ....
func (r *Status) Persist(tenantID string, status *models.Status) (*models.Status, error) {
	status.TenantID = tenantID
	if err := r.conn.GetConnectionWithPreload([]string{"SkillGroups"}).Create(status).Error; err != nil {
		return nil, err
	}
	return status, nil
}

// Update ....
func (r *Status) Update(tenantID string, status *models.Status) (*models.Status, error) {
	status.TenantID = tenantID
	if err := r.conn.GetConnectionWithPreload([]string{"SkillGroups"}).Save(status).Error; err != nil {
		return nil, err
	}
	return status, nil
}

// Delete ....
func (r *Status) Delete(tenantID string, status *models.Status) error {
	status.TenantID = tenantID
	if err := r.conn.GetConnection().Delete(status).Error; err != nil {
		return err
	}
	return nil
}

func (r *Status) GetByEmail(tenantID, email string) (*models.Status, error) {
	var status models.Status
	if err := r.conn.GetConnectionWithPreload(nil).Where("tenant_id = ?", tenantID).First(&status, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

// Get ....
func (r *Status) Get(tenantID string, ID string) (*models.Status, error) {
	var status models.Status
	if err := r.conn.GetConnectionWithPreload([]string{"SkillGroups"}).Where("tenant_id = ?", tenantID).First(&status, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

// GetAll ...
func (r *Status) GetAll(tenantID string, query database.Query, with ...string) ([]*models.Status, error) {
	var records []*models.Status
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

	tx.Model(&models.Status{}).
		Distinct().
		Where("tenant_id = ?", tenantID)

	if err := tx.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
