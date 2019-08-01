package load

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"hl/model"
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
		tx := db.Begin()
		for _, v := range users.Users {
			tx.Create(&v)
		}
		tx.Commit()
		//d := model.GetDate()

		//sql := "INSERT INTO User (id, email, first_name, last_name, gender, birth_date, age) VALUES "
		//for _, v := range users.Users {
		//	row := fmt.Sprintf("(%v), ", v)
		//	sql = sql + row
		//}
		//sql = sql[:len(sql)-2]
		//db.Exec(sql)
		fmt.Println(path)
		n++
		path = dir + "/users_" + strconv.Itoa(n) + ".json"
	}
	fmt.Println("finish")
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
		//sql := "INSERT INTO Location (id, distance, city, place, country) VALUES "
		tx := db.Begin()
		for _, v := range locs.Locations {
			tx.Create(&v)
		}
		tx.Commit()
		//sql = sql[:len(sql)-2]
		//db.Exec(sql)

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
		//sql := "INSERT INTO Visit (id, location, user, visited_at, mark) VALUES "
		tx := db.Begin()
		for _, v := range visits.Visits {
			tx.Create(&v)
		}
		tx.Commit()
		//sql = sql[:len(sql)-2]
		//db.Exec(sql)
		fmt.Println(path)
		n++
		path = dir + "/visits_" + strconv.Itoa(n) + ".json"
	}
}
