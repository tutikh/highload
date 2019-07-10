package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"highload/highload/model"
	"net/http"
	"strconv"
	"unicode"
)

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
func isLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func CreateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		RespondError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		RespondError(w, http.StatusInternalServerError)
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
		RespondError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := db.Save(&user).Error; err != nil {
		RespondError(w, http.StatusInternalServerError)
		return
	}
	return
}

func GetUserVisits(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		RespondError(w, http.StatusNotFound)
		return
	}

	query := db.Debug().Table("Visit").Select("Visit.mark, Visit.visited_at, Location.place").Joins("right join Location on Location.id = Visit.location").
		Where("Visit.user = ?", id)

	fromdate := r.FormValue("fromDate")
	if fromdate != "" {
		query = query.Where("Visit.visited_at > ?", fromdate)
	}
	_, ok := r.URL.Query()["fromDate"]
	if (ok && len(fromdate) < 1) || !isInt(fromdate) {
		RespondError(w, http.StatusBadRequest)
		return
	}

	todate := r.FormValue("toDate")
	if todate != "" {
		query = query.Where("Visit.visited_at < ?", todate)
	}
	_, ok = r.URL.Query()["toDate"]
	if (ok && len(todate) < 1) || !isInt(todate) {
		RespondError(w, http.StatusBadRequest)
		return
	}

	todistance := r.FormValue("toDistance")
	if todistance != "" {
		query = query.Where("Location.distance < ?", todistance)
	}
	_, ok = r.URL.Query()["toDistance"]
	if (ok && len(todistance) < 1) || !isInt(todistance) {
		RespondError(w, http.StatusBadRequest)
		return
	}

	country := r.FormValue("country")
	if country != "" {
		query = query.Where("Location.country = ?", country)
	}
	_, ok = r.URL.Query()["country"]
	if (ok && len(country) < 1) || !isLetter(country) {
		RespondError(w, http.StatusBadRequest)
		return
	}

	var result model.UserVisitsArray
	query.Order("Visit.visited_at").Scan(&result.Visits)
	respondJSON(w, http.StatusOK, result)
}

func getUserOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, model.User{ID: id}).Error; err != nil {
		RespondError(w, http.StatusNotFound)
		return nil
	}
	return &user
}
