package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Settings struct {
	Host     string
	Username string
	Password string
	Database string
	Client   *gorm.DB
}

func (db *Settings) InitializeDB(logging *log.Logger) (err error) {
	settings := make(map[string]string)
	settings["host"] = db.Host
	settings["user"] = db.Username
	settings["password"] = db.Password
	settings["database"] = db.Database
	settings["port"] = "5432"

	var dsn []string
	for key, value := range settings {
		dsn = append(dsn, key+"="+value)
	}

	fmt.Printf("Connecting to database: %s", strings.Join(dsn, " "))

	db.Client, err = gorm.Open(postgres.Open(strings.Join(dsn, " ")), &gorm.Config{Logger: logger.New(
		logging,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)})
	if err != nil {
		logging.Println(err)
		return err
	}

	tx := db.Client.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if tx.Error != nil {
		logging.Println(err)
		return err
	}

	return db.Client.AutoMigrate(Request{})
}
