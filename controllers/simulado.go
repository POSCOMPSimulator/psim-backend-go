package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetSimulados(w http.ResponseWriter, r *http.Request) {

	ok, user := utils.AuthUser(a.DB, w, r, 0)
	if !ok {
		return
	}

	var bsim models.BatchSimulados
	bsim.IDUsuario = user.Email

	if err := bsim.Get(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, bsim)

}

func (a *App) CreateSimulado(w http.ResponseWriter, r *http.Request) {

	ok, user := utils.AuthUser(a.DB, w, r, 0)
	if !ok {
		return
	}

	var sim models.Simulado
	sim.IdUsuario = user.Email
	sim.Estado = 0

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

func (a *App) GetSimulado(w http.ResponseWriter, r *http.Request) {

	ok, user := utils.AuthUser(a.DB, w, r, 0)
	if !ok {
		return
	}

	var err error
	var sim models.Simulado
	sim.IdUsuario = user.Email

	vars := mux.Vars(r)
	if id, ok := vars["id"]; ok {
		sim.ID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	if err := sim.Get(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, sim)

}

func (a *App) UpdateStateSimulado(w http.ResponseWriter, r *http.Request) {

	ok, user := utils.AuthUser(a.DB, w, r, 0)
	if !ok {
		return
	}

	var err error
	var sim models.Simulado
	sim.IdUsuario = user.Email

	vars := mux.Vars(r)
	if id, ok := vars["id"]; ok {
		sim.ID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	to_state, ok := vars["to_state"]
	if ok {
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Estado mal formatado.")
			return
		}
	}

	switch strings.ToUpper(to_state) {
	case "INICIAR":

		if err := sim.Start(a.DB); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		var retsim models.Simulado
		retsim.Questoes = sim.Questoes
		retsim.TempoRestante = sim.TempoRestante
		retsim.Respostas = sim.Respostas

		utils.RespondWithJSON(w, http.StatusAccepted, retsim)
		return

	case "CONTINUAR":

		if err := sim.Continue(a.DB); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		var retsim models.Simulado
		retsim.Questoes = sim.Questoes
		retsim.TempoRestante = sim.TempoRestante
		retsim.Respostas = sim.Respostas

		utils.RespondWithJSON(w, http.StatusAccepted, retsim)
		return

	case "FINALIZAR":

		defer r.Body.Close()

		if r.Body == http.NoBody {
			utils.RespondWithError(w, http.StatusBadRequest, "Body não encontrado.")
			return
		}

		var bresp models.BatchRespostas
		bresp.IDSimulado = sim.ID
		bresp.IDUsuario = user.Email

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&bresp); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := bresp.Update(a.DB); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err = sim.Finish(a.DB); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.WriteHeader(http.StatusMovedPermanently)

	default:

		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}

}

func (a *App) UpdateRespostasSimulado(w http.ResponseWriter, r *http.Request) {

	ok, user := utils.AuthUser(a.DB, w, r, 0)
	if !ok {
		return
	}

	var err error
	var bres models.BatchRespostas
	bres.IDUsuario = user.Email

	vars := mux.Vars(r)
	if id, ok := vars["id"]; ok {
		bres.IDSimulado, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	defer r.Body.Close()

	if r.Body == http.NoBody {
		utils.RespondWithError(w, http.StatusBadRequest, "Body não encontrado.")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&bres); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := bres.Update(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (a *App) DeleteSimulado(w http.ResponseWriter, r *http.Request) {

	ok, user := utils.AuthUser(a.DB, w, r, 0)
	if !ok {
		return
	}

	var err error
	var sim models.Simulado
	sim.IdUsuario = user.Email

	vars := mux.Vars(r)
	if id, ok := vars["id"]; ok {
		sim.ID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	sim.Delete(a.DB)
	w.WriteHeader(http.StatusOK)

}
