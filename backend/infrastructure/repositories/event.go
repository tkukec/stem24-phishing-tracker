package repositories

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
)

// NewEvent constructor for Event
func NewEvent(conn *database.Connection) EventRepository {
	return &Event{
		conn: conn,
	}
}

type EventRepository interface {
	Persist(tenantID string, event *models.Event) (*models.Event, error)
	Update(tenantID string, event *models.Event) (*models.Event, error)
	Delete(tenantID string, event *models.Event) error
	GetByEmail(tenantID, email string) (*models.Event, error)
	Get(tenantID string, ID string) (*models.Event, error)
	GetAll(tenantID string, query database.Query, with ...string) ([]*models.Event, error)
}

// Event ....
type Event struct {
	conn *database.Connection
}

// Persist ....
func (r *Event) Persist(tenantID string, event *models.Event) (*models.Event, error) {
	event.TenantID = tenantID
	if err := r.conn.GetConnectionWithPreload([]string{}).Create(event).Error; err != nil {
		return nil, err
	}
	return event, nil
}

// Update ....
func (r *Event) Update(tenantID string, event *models.Event) (*models.Event, error) {
	event.TenantID = tenantID
	if err := r.conn.GetConnectionWithPreload([]string{}).Save(event).Error; err != nil {
		return nil, err
	}
	return event, nil
}

// Delete ....
func (r *Event) Delete(tenantID string, event *models.Event) error {
	event.TenantID = tenantID
	if err := r.conn.GetConnection().Delete(event).Error; err != nil {
		return err
	}
	return nil
}

func (r *Event) GetByEmail(tenantID, email string) (*models.Event, error) {
	var event models.Event
	if err := r.conn.GetConnectionWithPreload(nil).Where("tenant_id = ?", tenantID).First(&event, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// Get ....
func (r *Event) Get(tenantID string, ID string) (*models.Event, error) {
	var event models.Event
	if err := r.conn.GetConnectionWithPreload([]string{}).Where("tenant_id = ?", tenantID).First(&event, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// GetAll ...
func (r *Event) GetAll(tenantID string, query database.Query, with ...string) ([]*models.Event, error) {
	var records []*models.Event
	preload := []string{"Comments", "Status"}
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

	tx.Model(&models.Event{}).
		Distinct().
		Where("tenant_id = ?", tenantID)

	if err := tx.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
