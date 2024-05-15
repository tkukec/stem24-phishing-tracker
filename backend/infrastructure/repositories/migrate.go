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
	if err = conn.GetConnection().AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err = conn.GetConnection().AutoMigrate(&models.Status{}); err != nil {
		return err
	}
	if err = conn.GetConnection().AutoMigrate(&models.Comment{}); err != nil {
		return err
	}
	if err = conn.GetConnection().AutoMigrate(&models.Event{}); err != nil {
		return err
	}
	return err
}

func DropTables(conn *database.Connection) {
	if err := conn.GetConnection().Migrator().DropTable(&models.Tenant{}); err != nil {
		log.Panic(err)
	}
	if err := conn.GetConnection().Migrator().DropTable(&models.User{}); err != nil {
		log.Panic(err)
	}
	if err := conn.GetConnection().Migrator().DropTable(&models.Status{}); err != nil {
		log.Panic(err)
	}
	if err := conn.GetConnection().Migrator().DropTable(&models.Comment{}); err != nil {
		log.Panic(err)
	}
	if err := conn.GetConnection().Migrator().DropTable(&models.Event{}); err != nil {
		log.Panic(err)
	}
}
