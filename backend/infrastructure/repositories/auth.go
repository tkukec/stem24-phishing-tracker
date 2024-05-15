package repositories

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
)

// NewUser constructor for User
func NewUser(conn *database.Connection) UserRepository {
	return &User{
		conn: conn,
	}
}

type UserRepository interface {
	Persist(tenantID string, user *models.User) (*models.User, error)
	Update(tenantID string, user *models.User) (*models.User, error)
	Delete(tenantID string, user *models.User) error
	GetByEmail(tenantID, email string) (*models.User, error)
	Get(tenantID string, ID string) (*models.User, error)
	GetAll(tenantID string, query database.Query, with ...string) ([]*models.User, error)
}

// User ....
type User struct {
	conn *database.Connection
}

// Persist ....
func (r *User) Persist(tenantID string, user *models.User) (*models.User, error) {
	user.TenantID = tenantID
	if err := r.conn.GetConnectionWithPreload([]string{"SkillGroups"}).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Update ....
func (r *User) Update(tenantID string, user *models.User) (*models.User, error) {
	user.TenantID = tenantID
	if err := r.conn.GetConnectionWithPreload([]string{"SkillGroups"}).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Delete ....
func (r *User) Delete(tenantID string, user *models.User) error {
	user.TenantID = tenantID
	if err := r.conn.GetConnection().Delete(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *User) GetByEmail(tenantID, email string) (*models.User, error) {
	var user models.User
	if err := r.conn.GetConnectionWithPreload(nil).Where("tenant_id = ?", tenantID).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Get ....
func (r *User) Get(tenantID string, ID string) (*models.User, error) {
	var user models.User
	if err := r.conn.GetConnectionWithPreload([]string{"SkillGroups"}).Where("tenant_id = ?", tenantID).First(&user, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll ...
func (r *User) GetAll(tenantID string, query database.Query, with ...string) ([]*models.User, error) {
	var records []*models.User
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

	tx.Model(&models.User{}).
		Distinct().
		Where("tenant_id = ?", tenantID)

	if err := tx.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
