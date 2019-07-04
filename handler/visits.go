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

func CreateVisit(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vis := model.Visit{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vis); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&vis).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, vis)
}

func GetVisit(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	vis := getVisitOr404(db, id, w, r)
	if vis == nil {
		return
	}
	respondJSON(w, http.StatusOK, vis)
}

func UpdateVisit(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	vis := getVisitOr404(db, id, w, r)
	if vis == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vis); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&vis).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, vis)
}

func getVisitOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Visit {
	vis := model.Visit{}
	if err := db.First(&vis, model.Visit{ID: id}).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &vis
}
