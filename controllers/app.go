package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/teris-io/shortid"
)

type App struct {
	Router    *gin.Engine
	DB        *sql.DB
	Cors      *cors.Cors
	AdminCode string
	SimIDGen  *shortid.Shortid
}

func (a *App) Initialize() error {

	_ = godotenv.Load()
	sid, _ := shortid.New(1, shortid.DefaultABC, 2342)

	var err error
	a.DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL")+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return err
	}

	a.Cors = cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowedOrigins:   []string{"*"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	a.AdminCode = os.Getenv("ADMIN_CODE")
	a.SimIDGen = sid

	a.Router = gin.Default()
	a.initializeRoutes()

	return nil

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Cors.Handler(a.Router)))
}
