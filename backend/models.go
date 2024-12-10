package backend

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Role struct {
	Id   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type User struct {
	Id        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	RoleId    uint      `json:"roleId" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Role      Role      `gorm:"references:id"`
}

type Session struct {
	Guid         string `json:"guid" gorm:"primaryKey"`
	RefreshToken string `json:"refresh_token" gorm:"unique;not null"`
	Ip           string `json:"ip" gorm:"not null"`
}

func DatabaseConnect() error {
	dsn := "host=localhost user=postgres password=postgres dbname=testcase port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	db.AutoMigrate(&Role{}, &User{})

	DB = db

	return nil
}
