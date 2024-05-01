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
)

type App struct {
	Router    *gin.Engine
	DB        *sql.DB
	Cors      *cors.Cors
	AdminCode string
}

func (a *App) Initialize() error {

	_ = godotenv.Load()

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

	a.Router = gin.Default()
	a.initializeRoutes()

	return nil

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Cors.Handler(a.Router)))
}
