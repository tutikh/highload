package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"highload/hl/model"
	"io/ioutil"
	"net/http"
	"strconv"
)

//var mu = &sync.Mutex{}
func CreateVisit(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//mu := &sync.Mutex{}
	//mu.Lock()
	//defer mu.Unlock()
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("pragma busy_timeout=5000;")
	//db.Exec("PRAGMA synchronous=normal;")
	//db.Exec("PRAGMA locking_mode=normal;")
	vis := model.Visit{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&vis); err != nil {
		RespondError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := db.First(&vis, model.Visit{ID: vis.ID}).Error; err == nil {
		RespondError(w, http.StatusBadRequest)
		return
	}
	if vis.ID != 0 && vis.User != 0 && vis.Location != 0 && vis.Mark != 0 && vis.VisitedAt != 0 {
		UpdateChan <- func() {
			db.Save(&vis)
		}
	} else {
		RespondError(w, http.StatusBadRequest)
		return
	}
	RespondJSON2(w, http.StatusOK)
}

func GetVisit(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//db.Exec("pragma busy_timeout=30000;")
	//db.Exec("PRAGMA journal_mode=DELETE;")
	//db.Exec("PRAGMA locking_mode=normal;")
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
	//mu := &sync.Mutex{}
	//mu.Lock()
	//defer mu.Unlock()
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("pragma busy_timeout=5000;")
	//db.Exec("PRAGMA synchronous=normal;")
	//db.Exec("PRAGMA locking_mode=normal;")
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

	req, _ := ioutil.ReadAll(r.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(req, &result); err != nil {
		RespondError(w, http.StatusBadRequest)
		return
	}
	for _, v := range result {
		if v == nil {
			RespondError(w, http.StatusBadRequest)
			return
		}
	}
	query := db.Model(vis)
	query.Updates(result)

	//if result["location"] != nil {
	//	query.Update("location", result["location"])
	//}
	//if result["user"] != nil {
	//	query.Update("user", result["user"])
	//}
	//if result["visited_at"] != nil {
	//	query.Update("visited_at", result["visited_at"])
	//}
	//if result["mark"] != nil {
	//	query.Update("mark", result["mark"])
	//}
	RespondJSON2(w, http.StatusOK)
}

func getVisitOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Visit {
	vis := model.Visit{}
	if err := db.First(&vis, model.Visit{ID: id}).Error; err != nil {
		RespondError(w, http.StatusNotFound)
		return nil
	}
	return &vis
}
