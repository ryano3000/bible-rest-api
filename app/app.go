package app

import (
	"bible/app/handler"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

type App struct {
	Environment string
	Router      *mux.Router
	DB          *gorm.DB
}

func (a *App) Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	a.Environment = os.Getenv("ENVIRONMENT")

	db, err := gorm.Open("sqlite3", "bible-sqlite.db")
	if err != nil {
		panic("Failed to connect database")
	}

	a.DB = db
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {
	a.Get("/translations", a.handleRequest(handler.GetAllTranslations))
	a.Get("/translations/{translation_id:[0-9]}/books", a.handleRequest(handler.GetAllBooks))
}

func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}
