package database

import (
	"auth/app/config"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(sqlite.Open(config.DB), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().Local() },
		Logger:  logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println("database connection error")
		panic(err)
	}
	DB = db
	fmt.Println("database connected")
}

func Migrate(tables ...interface{}) error {
	return DB.AutoMigrate(tables...)
}
