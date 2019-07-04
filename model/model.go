package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID         string `gorm:"unique" json:"id"`
	Email      string `json:"email"`
	First_name string `json:"firstname"`
	Last_name  string `json:"lastname"`
	Gender     string `json:"gender"`
	Birth_date int    `json:"birthdate"`
	Age        int    `json:"age"`
}

type Location struct {
	ID       uint
	Distance int
	City     string
	Place    string
	Country  string
}

type Visit struct {
	ID         uint
	Location   int
	User       int
	Visited_at int
	Mark       int
}

func (User) TableName() string {
	return "User"
}

func (Location) TableName() string {
	return "Location"
}

func (Visit) TableName() string {
	return "Visit"
}

func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&User{})
	return db
}
