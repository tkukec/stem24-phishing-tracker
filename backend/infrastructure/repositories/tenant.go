package repositories

import (
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
)

func NewTenant(conn *database.Connection) TenantRepository {
	return &Tenant{
		conn: conn,
	}
}

type TenantRepository interface {
	Persist(tenant *models.Tenant) (*models.Tenant, error)
	Update(Tenant *models.Tenant) (*models.Tenant, error)
	Delete(Tenant *models.Tenant) error
	GetByName(name string) (*models.Tenant, error)
	Get(ID string) (*models.Tenant, error)
	GetAll() ([]*models.Tenant, error)
}

type Tenant struct {
	conn *database.Connection
}

// Persist ....
func (r *Tenant) Persist(tenant *models.Tenant) (*models.Tenant, error) {
	if err := r.conn.GetConnection().Create(tenant).Error; err != nil {
		return nil, err
	}
	return tenant, nil
}

// Update ....
func (r *Tenant) Update(Tenant *models.Tenant) (*models.Tenant, error) {
	if err := r.conn.GetConnection().Save(Tenant).Error; err != nil {
		return nil, err
	}
	return Tenant, nil
}

// Delete ....
func (r *Tenant) Delete(Tenant *models.Tenant) error {
	if err := r.conn.GetConnection().Delete(Tenant).Error; err != nil {
		return err
	}
	return nil
}

func (r *Tenant) GetByName(name string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.conn.GetConnection().First(&tenant, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// Get ....
func (r *Tenant) Get(ID string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.conn.GetConnection().First(&tenant, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetAll ...
func (r *Tenant) GetAll() ([]*models.Tenant, error) {
	var records []*models.Tenant
	if err := r.conn.GetConnection().Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
