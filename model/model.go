package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	ID        int    `gorm:"unique" json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	BirthDate int    `json:"birth_date"`
	Age       int    `json:"age"`
}

type Location struct {
	ID       int    `gorm:"unique" json:"id"`
	Distance int    `json:"distance"`
	City     string `json:"city"`
	Place    string `json:"place"`
	Country  string `json:"country"`
}

type Visit struct {
	ID        int `gorm:"unique" json:"id"`
	Location  int `json:"location"`
	User      int `json:"user"`
	VisitedAt int `json:"visited_at"`
	Mark      int `json:"mark"`
}

type UserVisits struct {
	Mark      int    `json:"mark"`
	VisitedAt int    `json:"visited_at"`
	Place     string `json:"place"`
}

type UserVisitsArray struct {
	Visits []UserVisits `json:"visits"`
}

type LocationAvg struct {
	Avg float64 `json:"avg"`
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
