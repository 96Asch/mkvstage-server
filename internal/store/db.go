package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	RETRY_COUNT           = 10
	RETRY_TIME_IN_SECONDS = 5
)

// ConnectByName takes a DSN connection string and connects to a mysql database using GORM.
// Optional arguments can be supplied using args, which will add
// additional parameters to the DSN connection.
// It returns an error if the connection fails.
func GetDB(user, password, host, port, name string, args ...string) (*gorm.DB, error) {

	connectionString := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v",
		user,
		password,
		host,
		port,
		name,
	)

	if len(args) > 0 {
		connectionString += fmt.Sprintf("&%v", strings.Join(args, "&"))
	}

	ms, err := sql.Open("mysql", connectionString)
	log.Printf("Connecting to %v...", connectionString)

	if err != nil {
		return nil, err
	}

	for i := 0; i < RETRY_COUNT; i++ {
		err = ms.Ping()
		if err != nil {
			time.Sleep(RETRY_TIME_IN_SECONDS * time.Second)
			log.Printf("Retrying connection...")
			continue
		}

		database, err := gorm.Open(mysql.New(mysql.Config{
			Conn: ms,
		}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		log.Printf("Connection success")
		return database, err
	}

	return nil, errors.New("could not establish connection")

}
