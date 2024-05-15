package database

import (
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/logging"
	"gorm.io/gorm/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(driver, user, password, hostname, port, database string, debug bool, writer logger.Writer) (*Connection, error) {
	loggerToUse := logger.Default.LogMode(logger.Silent)
	if writer != nil {
		loggerToUse = logging.NewDbLogger(writer, logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Warn,
			Colorful:      false,
		}).LogMode(logger.Silent)
	}
	config := &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		Logger: loggerToUse,
	}

	var db *gorm.DB
	var err error
	switch driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", hostname, port, user, database, password)
		db, err = gorm.Open(postgres.Open(dsn), config)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database %s", err.Error())
		}
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, hostname, port, database)
		db, err = gorm.Open(mysql.Open(dsn), config)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database %s", err.Error())
		}
	default:
		dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, hostname, port, database)
		db, err = gorm.Open(mysql.Open(dsn), config)
		if err != nil {
			return nil, fmt.Errorf("unable to connect to database %s", err.Error())
		}
	}
	if debug {
		return NewConnectionWithDB(db.Debug()), nil
	}
	return NewConnectionWithDB(db), nil
}

// NewConnectionWithDB constructor for DatabaseConnection
func NewConnectionWithDB(db *gorm.DB) *Connection {
	return &Connection{
		db: db,
	}
}

// Connection ....
type Connection struct {
	db *gorm.DB
}

func (r *Connection) ConfigConnectionPooling(MaxIdleConns int, MaxOpenConns int, ConnMaxLifetime int) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("unable to get sql.DB connection pool interface: %s", err.Error())
	}

	if MaxIdleConns != 0 {
		sqlDB.SetMaxIdleConns(MaxIdleConns)
	}
	if MaxOpenConns != 0 {
		sqlDB.SetMaxOpenConns(MaxOpenConns)
	}
	if ConnMaxLifetime != 0 {
		sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(ConnMaxLifetime))
	}

	return nil
}

// GetConnection returns new gorm.DB connection.
func (r *Connection) GetConnection() *gorm.DB {
	//return r.db.Session(&gorm.Session{FullSaveAssociations: true})
	return r.db
}

func (r *Connection) TenantQueryConnection(tenantID string, with ...string) *gorm.DB {
	//return r.db.Session(&gorm.Session{FullSaveAssociations: true})
	if with == nil {
		with = make([]string, 0)
	}
	db := r.GetConnection().Where("tenant_id = ?", tenantID)
	for _, w := range with {
		db = db.Preload(w)
	}
	return db
}

// AddWith append preloads to the query
func (r *Connection) AddWith(db *gorm.DB, with []string) *gorm.DB {
	for _, w := range with {
		db = db.Preload(w)
	}
	return db
}

// GetConnectionWithPreload get db connection with preloads
func (r *Connection) GetConnectionWithPreload(with []string) *gorm.DB {
	if with == nil {
		with = make([]string, 0)
	}
	db := r.GetConnection()
	for _, w := range with {
		db = db.Preload(w)
	}
	return db
}
