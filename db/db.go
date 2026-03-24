package db

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Sighting struct {
	ID        uint   `gorm:"primaryKey"`
	Animal    string `gorm:"not null"`
	Location  string `gorm:"not null"`
	Notes     string
	CreatedAt int64 `gorm:"autoCreateTime"`
}

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("wildlife.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&Sighting{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
