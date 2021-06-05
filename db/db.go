package db

import (
	"os"

	"github.com/morsby/kbu"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	if _, err := os.Stat("db.sqlite"); err != os.ErrNotExist {
		os.Remove("db.sqlite")
	}

	return gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{
		CreateBatchSize: 1000,
	})
}

func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}

func InsertRounds(db *gorm.DB, rounds *[]kbu.Round) error {
	if res := db.Create(rounds); res.Error != nil {
		return res.Error
	}
	return nil
}
