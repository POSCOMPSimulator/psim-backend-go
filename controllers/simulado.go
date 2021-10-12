package controllers

import (
	"encoding/json"
	"net/http"

	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetSimulados(w http.ResponseWriter, r *http.Request) {}

func (a *App) CreateSimulado(w http.ResponseWriter, r *http.Request) {

	ok, user := utils.AuthUser(a.DB, w, r, 1)
	if !ok {
		return
	}

	var sim models.Simulado
	sim.IdUsuario = user.Email

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sim); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := sim.Create(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (a *App) GetSimulado(w http.ResponseWriter, r *http.Request) {}

func (a *App) UpdateStateSimulado(w http.ResponseWriter, r *http.Request) {}

func (a *App) UpdateRespostasSimulado(w http.ResponseWriter, r *http.Request) {}

func (a *App) DeleteSimulado(w http.ResponseWriter, r *http.Request) {}
