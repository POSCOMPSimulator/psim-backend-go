package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/models"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize() {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", goDotEnvVariable("APP_DB_USERNAME"), goDotEnvVariable("APP_DB_PASSWORD"), goDotEnvVariable("APP_DB_NAME"))

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", (getCorsMiddleware()).Handler(a.Router)))
}

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/questao/", a.getQuestoes).Methods("GET")
	a.Router.HandleFunc("/questao/", a.createQuestao).Methods("POST")
	a.Router.HandleFunc("/questao/sumario/", a.getQSumario).Methods("GET")
	a.Router.HandleFunc("/questao/{id}/", a.getQuestao).Methods("GET")
	a.Router.HandleFunc("/questao/{id}/", a.reportQuestao).Methods("PUT")
	a.Router.HandleFunc("/questao/{id}/", a.updateQuestao).Methods("PATCH")
	a.Router.HandleFunc("/questao/{id}/", a.deleteQuestao).Methods("DELETE")
}

func (a *App) handleFetch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching")
}

func (a *App) getQuestoes(w http.ResponseWriter, r *http.Request) {}

func (a *App) getQSumario(w http.ResponseWriter, r *http.Request) {}

func (a *App) getQuestao(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID inválido.")
		return
	}

	q := models.Questao{ID: id}
	if err := q.GetQuestao(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, q)

}

func (a *App) createQuestao(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		return
	}

	_, err := auth.VerifyIdToken(r.Header["User-Token"][0])

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var q models.Questao
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&q); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := q.CreateQuestao(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, q)
}

func (a *App) reportQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) updateQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) deleteQuestao(w http.ResponseWriter, r *http.Request) {}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func getCorsMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "FETCH", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "user-token", "User-Token"},
		// Enable Debugging for testing, consider disabling in production
		Debug:              true,
		OptionsPassthrough: true,
	})
}
