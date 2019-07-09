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

func CreateLocation(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	loc := model.Location{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loc); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&loc).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	return
}

func GetLocation(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	loc := getLocationOr404(db, id, w, r)
	if loc == nil {
		return
	}
	respondJSON(w, http.StatusOK, loc)
}

func UpdateLocation(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	loc := getLocationOr404(db, id, w, r)
	if loc == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loc); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&loc).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	return
}

func GetAvg(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
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
	fromage := r.FormValue("fromAge")
	toage := r.FormValue("toAge")
	if toage == "" {
		toage = "100"
	}
	gender := r.FormValue("gender")

	query := db.Debug().Table("Visit").Select("ROUND(AVG(Visit.mark), 5) as avg").Joins("right join User on User.id = Visit.user").
		Where("Visit.location = ? AND Visit.visited_at > ? AND Visit.visited_at < ? AND User.age > ? AND User.age < ?", id, fromdate, todate, fromage, toage)

	if gender != "" {
		query = db.Debug().Table("Visit").Select("ROUND(AVG(Visit.mark), 5) as avg").Joins("right join User on User.id = Visit.user").
			Where("Visit.location = ? AND Visit.visited_at > ? AND Visit.visited_at < ? AND User.age > ? AND User.age < ? AND User.gender = ?", id, fromdate, todate, fromage, toage, gender)
	}

	var result model.LocationAvg
	query.Scan(&result)

	if result.Avg == float64(int64(result.Avg)) {
		respondJSONforInt(w, http.StatusOK, result.Avg)
	} else {
		respondJSON(w, http.StatusOK, result)
	}
}

func getLocationOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Location {
	loc := model.Location{}
	if err := db.First(&loc, model.Location{ID: id}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &loc
}
