package load

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"highload/hl/model"
	"io/ioutil"
	"os"
	"strconv"
)

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func LoadUser(db *gorm.DB, dir string) {
	db.Delete(&model.User{})
	n := 1
	path := dir + "/users_" + strconv.Itoa(n) + ".json"

	for Exists(path) {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("cant open file %v", err.Error())
			os.Exit(100)
		}
		defer file.Close()
		read, _ := ioutil.ReadAll(file)
		var users model.Users
		if err := json.Unmarshal(read, &users); err != nil {
			fmt.Printf("cant unmarshal %v", err.Error())
			os.Exit(100)
		}

		d := model.GetDate()

		sql := "INSERT INTO User (id, email, first_name, last_name, gender, birth_date, age) VALUES "
		for _, v := range users.Users {
			row := fmt.Sprintf("(%d, '%s', '%s', '%s', '%s', %d, %d), ", v.ID, v.Email, v.FirstName, v.LastName, v.Gender, v.BirthDate, (d-v.BirthDate)/31536000)
			sql = sql + row
		}
		sql = sql[:len(sql)-2]
		db.Exec(sql)
		fmt.Println(path)
		n++
		path = dir + "/users_" + strconv.Itoa(n) + ".json"
	}
}

func LoadLocation(db *gorm.DB, dir string) {
	db.Delete(&model.Location{})
	n := 1
	path := dir + "/locations_" + strconv.Itoa(n) + ".json"

	for Exists(path) {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("cant open file %v", err.Error())
			os.Exit(100)
		}
		defer file.Close()
		read, _ := ioutil.ReadAll(file)
		var locs model.Locations
		if err := json.Unmarshal(read, &locs); err != nil {
			fmt.Printf("cant unmarshal %v", err.Error())
			os.Exit(100)
		}
		sql := "INSERT INTO Location (id, distance, city, place, country) VALUES "
		for _, v := range locs.Locations {
			row := fmt.Sprintf("(%d, %d, '%s', '%s', '%s'), ", v.ID, v.Distance, v.City, v.Place, v.Country)
			sql = sql + row
		}
		sql = sql[:len(sql)-2]
		db.Exec(sql)
		fmt.Println(path)
		n++
		path = dir + "/locations_" + strconv.Itoa(n) + ".json"
	}
}

func LoadVisit(db *gorm.DB, dir string) {
	db.Delete(&model.Visit{})
	n := 1
	path := dir + "/visits_" + strconv.Itoa(n) + ".json"

	for Exists(path) {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("cant open file %v", err.Error())
			os.Exit(100)
		}
		defer file.Close()
		read, _ := ioutil.ReadAll(file)
		var visits model.Visits
		if err := json.Unmarshal(read, &visits); err != nil {
			fmt.Printf("cant unmarshal %v", err.Error())
			os.Exit(100)
		}
		sql := "INSERT INTO Visit (id, location, user, visited_at, mark) VALUES "
		for _, v := range visits.Visits {
			row := fmt.Sprintf("(%d, %d, %d, %d, %d), ", v.ID, v.Location, v.User, v.VisitedAt, v.Mark)
			sql = sql + row
		}
		sql = sql[:len(sql)-2]
		db.Exec(sql)
		fmt.Println(path)
		n++
		path = dir + "/visits_" + strconv.Itoa(n) + ".json"
	}
}
