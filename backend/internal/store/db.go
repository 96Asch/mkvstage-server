package store

import (
	"fmt"
	"log"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectByName takes a DSN connection string and connects to a mysql database using GORM.
// Optional arguments can be supplied using args, which will add
// additional parameters to the DSN connection.
// It returns an error if the connection fails.
func GetDB(user, password, host, port, name string) (*gorm.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, port, name,
	)

	database, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	log.Printf("Connecting to %v...", connectionString)

	if err != nil {
		return nil, domain.NewInitializationErr(err.Error())
	}

	log.Printf("Connection Succeeded")

	return database, nil
}
