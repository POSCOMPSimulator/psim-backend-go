package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/models"

	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetQuestoes(w http.ResponseWriter, r *http.Request) {}

func (a *App) GetQSumario(w http.ResponseWriter, r *http.Request) {}

func (a *App) GetQuestao(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "ID inválido.")
		return
	}

	q := models.Questao{ID: id}
	if err := q.GetQuestao(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, q)

}

func (a *App) CreateQuestao(w http.ResponseWriter, r *http.Request) {

	_, valid, err := auth.VerifyIdToken(r.Header["User-Token"][0])

	if !valid {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var q models.Questao
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&q); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := q.CreateQuestao(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, q)
}

func (a *App) ReportQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) UpdateQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) DeleteQuestao(w http.ResponseWriter, r *http.Request) {}
