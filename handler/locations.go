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
	respondJSON(w, http.StatusCreated, loc)
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
	respondJSON(w, http.StatusOK, loc)
}

func GetAvg(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	var result model.LocationAvg
	db.Table("Visit").Select("AVG(Visit.mark) as avg").Where("Visit.location = ?", id).Scan(&result)
	respondJSON(w, http.StatusOK, result)
}

func getLocationOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Location {
	loc := model.Location{}
	if err := db.First(&loc, model.Location{ID: id}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &loc
}
