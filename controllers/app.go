package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Cors   *cors.Cors
}

func (a *App) Initialize() {

	_ = godotenv.Load()

	var err error
	a.DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	a.Cors = cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Cors.Handler(a.Router)))
}
