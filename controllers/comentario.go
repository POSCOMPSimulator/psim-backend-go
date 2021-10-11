package controllers

import (
	"net/http"

	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetComentariosSinalizados(w http.ResponseWriter, r *http.Request) {

	if !utils.AuthUserModerator(a.DB, w, r, 1) {
		return
	}

	var bc models.BatchComentarios
	if err := bc.GetComentariosSinalizados(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	utils.RespondWithJSON(w, http.StatusOK, bc)

}

func (a *App) GetComentariosQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) PostComentarioQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) ReportComentario(w http.ResponseWriter, r *http.Request) {}

func (a *App) DeleteComentario(w http.ResponseWriter, r *http.Request) {}
