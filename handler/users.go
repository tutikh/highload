package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"highload/highload/model"
	"net/http"
	"strconv"
)

func CreateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	return
}

func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	user := getUserOr404(db, id, w, r)
	if user == nil {
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func UpdateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	user := getUserOr404(db, id, w, r)
	if user == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	return
}

func GetUserVisits(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}

	fromdate := r.FormValue("fromDate")
	todate := r.FormValue("toDate")
	if todate == "" {
		todate = "9999999999"
	}
	todistance := r.FormValue("toDistance")
	if todistance == "" {
		todistance = "9999999999"
	}

	country := r.FormValue("country")

	query := db.Debug().Table("Visit").Select("Visit.mark, Visit.visited_at, Location.place").Joins("right join Location on Location.id = Visit.location").
		Where("Visit.user = ? AND Visit.visited_at > ? AND Visit.visited_at < ? AND Location.distance < ?", id, fromdate, todate, todistance)

	if country != "" {
		query = db.Debug().Table("Visit").Select("Visit.mark, Visit.visited_at, Location.place").Joins("right join Location on Location.id = Visit.location").
			Where("Visit.user = ? AND Visit.visited_at > ? AND Visit.visited_at < ? AND Location.distance < ? AND Location.country = ?", id, fromdate, todate, todistance, country)
	}

	var result model.UserVisitsArray
	query.Order("Visit.visited_at").Scan(&result.Visits)
	respondJSON(w, http.StatusOK, result)
}

func getUserOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, model.User{ID: id}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &user
}
