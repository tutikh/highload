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
	////mu := &sync.Mutex{}
	//
	//db.Exec("PRAGMA journal_mode=WAL;")
	//db.Exec("pragma busy_timeout=5000;")
	//db.Exec("PRAGMA synchronous=normal;")
	//db.Exec("PRAGMA locking_mode=EXCLUSIVE;")

	user := model.User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		//fmt.Println("Decoding problem(u)")
		RespondError(w, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	//if err := db.First(&user, model.User{ID: user.ID}).Error; err == nil {
	//	fmt.Println("User is exist: ", &user)
	//	RespondError(w, http.StatusBadRequest)
	//	return
	//}
	if user.ID == 0 || user.Email == "" || user.BirthDate == 0 || user.FirstName == "" || user.Gender == "" || user.LastName == "" {
		//fmt.Println("Bad request")
		RespondError(w, http.StatusBadRequest)
		return
	}
	//mu.Lock()
	//defer mu.Unlock()
	err := db.Save(&user).Error
	if err != nil {
		fmt.Println("not saved(u)")
	}
	RespondJSON2(w, http.StatusOK)
	//fmt.Println("saved(u)")
	return
}

func GetUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//db.Exec("pragma busy_timeout=30000;")
	//db.Exec("PRAGMA journal_mode=OFF;")
	//db.Exec("PRAGMA locking_mode=normal;")
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
	//mu := &sync.Mutex{}

	//db.Exec("PRAGMA journal_mode=WAL;")
	//db.Exec("pragma busy_timeout=5000;")
	//db.Exec("PRAGMA synchronous=normal;")
	vars := mux.Vars(r)
	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		fmt.Print(err)
	}
	user := getUserOr404(db, id, w, r)
	if user == nil {
		//fmt.Println("User not found on update")
		return
	}
	req, _ := ioutil.ReadAll(r.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(req, &result); err != nil {
		//fmt.Println("Unmarshaling problem")
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
	//mu.Lock()
	//defer mu.Unlock()

	if err := json.Unmarshal(req, &user); err != nil {
		//fmt.Println("Unmarshaling problem")
		RespondError(w, http.StatusBadRequest)
		return
	}
	db.Save(&user)
	//query.Updates(result)

	//if result["first_name"] != nil {
	//	query.Update("first_name", result["first_name"])
	//}
	//if result["last_name"] != nil {
	//	query.Update("last_name", result["last_name"])
	//}
	//if result["email"] != nil {
	//	query.Update("email", result["email"])
	//}
	//if result["gender"] != nil {
	//	query.Update("gender", result["gender"])
	//}
	//if result["birth_date"] != nil {
	//	query.Update("birth_date", result["birth_date"])
	//}
	//fmt.Println("saved")
	RespondJSON2(w, http.StatusOK)
	return
}

func GetUserVisits(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	//mu.Lock()
	//defer mu.Unlock()
	//db.Exec("pragma busy_timeout=30000;")
	//db.Exec("PRAGMA journal_mode=OFF;")
	//db.Exec("PRAGMA locking_mode=normal;")
	vars := mux.Vars(r)
	v := vars["id"]
	id, err := strconv.Atoi(v)
	if err != nil {
		RespondError(w, http.StatusNotFound)
		return
	}
	user := model.User{}
	if err := db.First(&user, model.User{ID: id}).Error; err != nil {
		RespondError(w, http.StatusNotFound)
		return
	}
	query := db.Table("Visit").Select("Visit.mark, Visit.visited_at, Location.place").Joins("inner join Location on Location.id = Visit.location").
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
	if ok && len(country) < 1 {
		RespondError(w, http.StatusBadRequest)
		return
	}

	var result model.UserVisitsArray
	query.Order("Visit.visited_at").Scan(&result.Visits)
	respondJSON(w, http.StatusOK, result)
}

func getUserOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.User {
	//mu.Lock()
	//defer mu.Unlock()
	user := model.User{}
	if err := db.First(&user, model.User{ID: id}).Error; err != nil {
		RespondError(w, http.StatusNotFound)
		return nil
	}
	return &user
}
