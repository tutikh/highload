package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"hl/model"
	"io/ioutil"
	"net/http"
	"strconv"
)

func CreateLocation(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//mu := &sync.Mutex{}

	//db.Exec("PRAGMA journal_mode=WAL;")
	//db.Exec("pragma busy_timeout=5000;")
	//db.Exec("PRAGMA synchronous=normal;")
	//db.Exec("PRAGMA locking_mode=EXCLUSIVE;")
	loc := model.Location{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loc); err != nil {
		//fmt.Println("Decoding problem(loc)")
		RespondError(w, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	//err := db.First(&loc, model.Location{ID: loc.ID}).Error
	//if err == nil {
	//	//fmt.Println("Location is exist")
	//	RespondError(w, http.StatusBadRequest)
	//	return
	//}
	//if err != gorm.ErrRecordNotFound {
	//	fmt.Println(err.Error())
	//}
	if loc.ID == 0 || loc.City == "" || loc.Country == "" || loc.Distance == 0 || loc.Place == "" {
		//fmt.Println("Bad request")
		RespondError(w, http.StatusBadRequest)
		return
	}
	//mu.Lock()
	//defer mu.Unlock()
	err := db.Save(&loc).Error
	if err != nil {
		fmt.Println("loc not saved")
	}
	RespondJSON2(w, http.StatusOK)
	//fmt.Println("loc saved")
	return
}

func GetLocation(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//db.Exec("pragma busy_timeout=30000;")
	//db.Exec("PRAGMA journal_mode=OFF;")
	//db.Exec("PRAGMA locking_mode=normal;")
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
	//mu := &sync.Mutex{}

	//db.Exec("PRAGMA journal_mode=WAL;")
	//db.Exec("pragma busy_timeout=5000;")
	//db.Exec("PRAGMA synchronous=normal;")
	//db.Exec("PRAGMA locking_mode=EXCLUSIVE;")
	vars := mux.Vars(r)

	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	loc := getLocationOr404(db, id, w, r)
	if loc == nil {
		//fmt.Println("Loaction not found(89)")
		return
	}
	req, _ := ioutil.ReadAll(r.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(req, &result); err != nil {
		//fmt.Println("Unmarshaling problem(loc)")
		RespondError(w, http.StatusBadRequest)
		return
	}
	for _, v := range result {
		if v == nil {
			//fmt.Println("problem")
			RespondError(w, http.StatusBadRequest)
			return
		}
	}
	query := db.Model(loc)
	//mu.Lock()
	//defer mu.Unlock()
	err = query.Updates(result).Error
	if err != nil {
		fmt.Println("loc not updated")
	}

	//if result["distance"] != nil {
	//	query.Update("distance", result["distance"])
	//}
	//if result["city"] != nil {
	//	query.Update("city", result["city"])
	//}
	//if result["place"] != nil {
	//	query.Update("place", result["place"])
	//}
	//if result["country"] != nil {
	//	query.Update("country", result["country"])
	//}
	RespondJSON2(w, http.StatusOK)
	return
}

func GetAvg(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//mu.Lock()
	//defer mu.Unlock()
	//db.Exec("pragma busy_timeout=30000;")
	//
	//db.Exec("PRAGMA journal_mode=OFF;")
	//db.Exec("PRAGMA locking_mode=normal;")
	vars := mux.Vars(r)
	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	loc := model.Location{}
	if err := db.First(&loc, model.Location{ID: id}).Error; err != nil {
		RespondError(w, http.StatusNotFound)
		return
	}
	query := db.Table("Visit").Select("ROUND(AVG(Visit.mark), 5) as avg").Joins("inner join User on User.id = Visit.user").
		Where("Visit.location = ?", id)

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

	fromage := r.FormValue("fromAge")
	if fromage != "" {
		query = query.Where("User.age >= ?", fromage)
	}
	_, ok = r.URL.Query()["fromAge"]
	if (ok && len(fromage) < 1) || !isInt(fromage) {
		RespondError(w, http.StatusBadRequest)
		return
	}

	toage := r.FormValue("toAge")
	if toage != "" {
		query = query.Where("User.age < ?", toage)
	}
	_, ok = r.URL.Query()["toAge"]
	if (ok && len(toage) < 1) || !isInt(toage) {
		RespondError(w, http.StatusBadRequest)
		return
	}

	gender := r.FormValue("gender")
	if gender != "" {
		query = query.Where("User.gender = ?", gender)
	}
	_, ok = r.URL.Query()["gender"]
	if (ok && len(gender) != 1) || !isLetter(gender) {
		RespondError(w, http.StatusBadRequest)
		return
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
	//mu.Lock()
	//defer mu.Unlock()
	loc := model.Location{}
	if err := db.First(&loc, model.Location{ID: id}).Error; err != nil {
		RespondError(w, http.StatusNotFound)
		return nil
	}
	return &loc
}
