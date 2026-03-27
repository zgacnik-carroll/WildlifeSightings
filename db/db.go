package db

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	CreatedAt int64  `gorm:"autoCreateTime"`
}

type Sighting struct {
	ID        uint   `gorm:"primaryKey"`
	Animal    string `gorm:"not null"`
	Location  string `gorm:"not null"`
	Notes     string
	UserID    uint  `gorm:"not null"`
	User      User  `gorm:"foreignKey:UserID"`
	CreatedAt int64 `gorm:"autoCreateTime"`
}

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("wildlife.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&User{}, &Sighting{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	dropLegacyEmailColumn()
}

func dropLegacyEmailColumn() {
	if !DB.Migrator().HasTable("users") {
		return
	}

	if !DB.Migrator().HasColumn("users", "email") {
		return
	}

	if DB.Migrator().HasIndex("users", "idx_users_email") {
		if err := DB.Migrator().DropIndex("users", "idx_users_email"); err != nil {
			log.Fatal("Failed to drop email index:", err)
		}
	}

	if err := DB.Exec("ALTER TABLE users DROP COLUMN email").Error; err != nil {
		log.Fatal("Failed to drop legacy email column:", err)
	}
}
