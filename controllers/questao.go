package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"poscomp-simulator.com/backend/models/questao"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetQuestoes(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var batch questao.BatchQuestoes
	if val, ok := r.Form["anos"]; ok {
		batch.Filtros.Anos = make([]int, len(val))
		for e, v := range val {
			i, err := strconv.Atoi(v)

			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Campo anos mal formatado")
				return
			}

			batch.Filtros.Anos[e] = i
		}
	}

	if val, ok := r.Form["areas"]; ok {
		batch.Filtros.Areas = val
	}

	if _, ok := r.Form["sinalizadas"]; ok {
		batch.Filtros.Sinalizadas = true
	}

	if err := batch.Get(a.DB); err != nil {
		fmt.Println(err)
	}

	utils.RespondWithJSON(w, http.StatusOK, batch)

}

func (a *App) GetQSumario(w http.ResponseWriter, r *http.Request) {

	var sq questao.SumarioQuestoes
	sq.Get(a.DB)
	utils.RespondWithJSON(w, http.StatusOK, sq)

}

func (a *App) GetErrosQuestao(w http.ResponseWriter, r *http.Request) {

	if !utils.AuthUserModerator(a.DB, w, r) {
		return
	}

	var errosq questao.ErrosQuestao
	var err error
	vars := mux.Vars(r)

	if value, ok := vars["id"]; ok {
		errosq.ID, err = strconv.Atoi(value)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	errosq.Get(a.DB)
	utils.RespondWithJSON(w, http.StatusOK, errosq)

}

func (a *App) SolveErrosQuestao(w http.ResponseWriter, r *http.Request) {

	if !utils.AuthUserModerator(a.DB, w, r) {
		return
	}

	var errosq questao.ErrosQuestao
	vars := mux.Vars(r)
	var err error
	if value, ok := vars["id"]; ok {
		errosq.ID, err = strconv.Atoi(value)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&errosq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	errosq.Solve(a.DB)

}

func (a *App) CreateQuestao(w http.ResponseWriter, r *http.Request) {

	if !utils.AuthUserModerator(a.DB, w, r) {
		return
	}

	var q questao.Questao
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&q); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := q.Create(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (a *App) ReportQuestao(w http.ResponseWriter, r *http.Request) {

	var m questao.MensagemErro

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := m.Report(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

}

func (a *App) UpdateQuestao(w http.ResponseWriter, r *http.Request) {

	if !utils.AuthUserModerator(a.DB, w, r) {
		return
	}

	var q questao.Questao
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&q); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := q.Update(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

}

func (a *App) DeleteQuestao(w http.ResponseWriter, r *http.Request) {

	if !utils.AuthUserModerator(a.DB, w, r) {
		return
	}

	vars := mux.Vars(r)
	var err error
	var q questao.Questao
	if id, ok := vars["id"]; ok {
		q.ID, err = strconv.Atoi(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusBadRequest, "ID mal formatado.")
		return
	}

	q.Delete(a.DB)

}
