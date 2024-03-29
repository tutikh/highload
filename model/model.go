package model

import (
	"bufio"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"hl/config"
	"os"
	"strconv"
)

type User struct {
	ID        int    `gorm:"unique" json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	BirthDate int    `json:"birth_date"`
	Age       int    `json:"-"`
}

type Users struct {
	Users []User `json:"users"`
}

type Location struct {
	ID       int    `gorm:"unique" json:"id"`
	Distance int    `json:"distance"`
	City     string `json:"city"`
	Place    string `json:"place"`
	Country  string `json:"country"`
}

type Locations struct {
	Locations []Location `json:"locations"`
}

type Visit struct {
	ID        int `gorm:"unique" json:"id"`
	Location  int `json:"location"`
	User      int `json:"user"`
	VisitedAt int `json:"visited_at"`
	Mark      int `json:"mark"`
}

type Visits struct {
	Visits []Visit `json:"visits"`
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
	db.AutoMigrate(&User{}, &Location{}, &Visit{})
	return db
}

func GetDate() int {
	var d int

	config := config.GetConfig("/root/go/src/hl/config/Config2.json")
	f, err := os.Open(config.Dataoptions)
	if err != nil {
		fmt.Printf("cant open file %v", err.Error())
		os.Exit(100)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		d, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Printf("%v", err.Error())
		}
		return d
	}
	return d
}

func (u *User) BeforeSave() {
	d := GetDate()
	u.Age = (d - u.BirthDate) / 31536000
	return
}

//func (u *User) BeforeUpdate() {
//	d := GetDate()
//	u.Age = (d - u.BirthDate) / 31536000
//	return
//}
