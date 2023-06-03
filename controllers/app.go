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
	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/mailer"
)

type App struct {
	Router     *gin.Engine
	DB         *sql.DB
	tokenMaker auth.Maker
	Cors       *cors.Cors
	Mailer     *mailer.Mailer
}

func (a *App) Initialize() error {

	_ = godotenv.Load()

	var err error
	a.DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
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

	a.tokenMaker, err = auth.NewPasetoMaker(os.Getenv("TOKEN_SYMMETRIC_KEY"))

	if err != nil {
		log.Fatal(err)
		return err
	}

	a.Mailer = mailer.NewMailer(
		os.Getenv("SENDER_EMAIL"),
		os.Getenv("SENDER_PASSWORD"),
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
	)

	a.Router = gin.Default()
	a.initializeRoutes()

	return nil

}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Cors.Handler(a.Router)))
}
