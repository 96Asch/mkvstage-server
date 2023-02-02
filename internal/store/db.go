package store

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectByName takes a DSN connection string and connects to a mysql database using GORM.
// Optional arguments can be supplied using args, which will add
// additional parameters to the DSN connection.
// It returns an error if the connection fails.
func GetDB(host, port, name, user, password string, args ...string) (*gorm.DB, error) {

	connectionString := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v&%v",
		host,
		port,
		name,
		user,
		password,
		strings.Join(args, "&"))

	database, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	return database, err
}
