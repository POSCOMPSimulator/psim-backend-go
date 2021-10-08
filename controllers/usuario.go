package controllers

import (
	"net/http"

	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) CreateOrLoginUsuario(w http.ResponseWriter, r *http.Request) {}

func (a *App) GetUsuario(w http.ResponseWriter, r *http.Request) {

	user, err := auth.VerifyIdToken(r.Header.Get("Authorization"))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err = user.Get(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
	return
}

func (a *App) PromoteUsuario(w http.ResponseWriter, r *http.Request) {}

func (a *App) DeleteUsuario(w http.ResponseWriter, r *http.Request) {}
