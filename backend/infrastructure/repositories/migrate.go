package repositories

import (
	"github.com/andrezz-b/stem24-phishing-tracker/domain/models"
	database "github.com/andrezz-b/stem24-phishing-tracker/shared/database"
	"log"
)

func Migrate(conn *database.Connection) error {

	var err error
	if err = conn.GetConnection().AutoMigrate(&models.Tenant{}); err != nil {
		return err
	}
	if err = conn.GetConnection().AutoMigrate(&models.Agent{}); err != nil {
		return err
	}
	return err
}

func DropTables(conn *database.Connection) {
	if err := conn.GetConnection().Migrator().DropTable(&models.Agent{}); err != nil {
		log.Panic(err)
	}
	if err := conn.GetConnection().Migrator().DropTable(&models.Tenant{}); err != nil {
		log.Panic(err)
	}
}
