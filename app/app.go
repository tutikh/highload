package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"highload/highload/config"
	"highload/highload/handler"
	"highload/highload/model"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.Router.NotFoundHandler = http.HandlerFunc(Custom404)
	a.setRouters()
}
func Custom404(w http.ResponseWriter, r *http.Request) {
	handler.RespondError(w, http.StatusNotFound)
	return
}

func (a *App) setRouters() {
	a.Post("/users/new", a.CreateUser)
	a.Get("/users/{id:[0-9]+}", a.GetUser)
	a.Post("/users/{id:[0-9]+}", a.UpdateUser)
	a.Post("/locations/new", a.CreateLocation)
	a.Get("/locations/{id:[0-9]+}", a.GetLocation)
	a.Post("/locations/{id:[0-9]+}", a.UpdateLocation)
	a.Post("/visits/new", a.CreateVisit)
	a.Get("/visits/{id:[0-9]+}", a.GetVisit)
	a.Post("/visits/{id:[0-9]+}", a.UpdateVisit)
	a.Get("/users/{id:[0-9]+}/visits", a.GetUserVisits)
	a.Get("/locations/{id:[0-9]+}/avg", a.GetAvg)
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	handler.CreateUser(a.DB, w, r)
}

func (a *App) GetUser(w http.ResponseWriter, r *http.Request) {
	handler.GetUser(a.DB, w, r)
}

func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	handler.UpdateUser(a.DB, w, r)
}

func (a *App) CreateLocation(w http.ResponseWriter, r *http.Request) {
	handler.CreateLocation(a.DB, w, r)
}

func (a *App) GetLocation(w http.ResponseWriter, r *http.Request) {
	handler.GetLocation(a.DB, w, r)
}

func (a *App) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	handler.UpdateLocation(a.DB, w, r)
}

func (a *App) CreateVisit(w http.ResponseWriter, r *http.Request) {
	handler.CreateVisit(a.DB, w, r)
}

func (a *App) GetVisit(w http.ResponseWriter, r *http.Request) {
	handler.GetVisit(a.DB, w, r)
}

func (a *App) UpdateVisit(w http.ResponseWriter, r *http.Request) {
	handler.UpdateVisit(a.DB, w, r)
}

func (a *App) GetUserVisits(w http.ResponseWriter, r *http.Request) {
	handler.GetUserVisits(a.DB, w, r)
}

func (a *App) GetAvg(w http.ResponseWriter, r *http.Request) {
	handler.GetAvg(a.DB, w, r)
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
