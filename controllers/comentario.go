package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetComentariosSinalizados(w http.ResponseWriter, r *http.Request) {

	if !utils.AuthUser(a.DB, w, r, 1) {
		return
	}

	var bc models.BatchComentarios
	if err := bc.GetComentariosSinalizados(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, bc)

}

func (a *App) GetComentariosQuestao(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var err error
	var bc models.BatchComentarios

	if id, ok := vars["id"]; ok {
		bc.QuestaoID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
	}

	if err := bc.GetComentariosQuestao(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, bc)

}

func (a *App) PostComentarioQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) ReportComentario(w http.ResponseWriter, r *http.Request) {}

func (a *App) DeleteComentario(w http.ResponseWriter, r *http.Request) {}
