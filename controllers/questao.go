package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetQuestoes(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var batch models.BatchQuestoes
	if val, ok := r.Form["anos"]; ok {
		batch.Filtros.Anos = make([]int, len(val))
		for e, v := range val {
			i, err := strconv.Atoi(v)

			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Campo anos mal formatado")
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

func (a *App) GetQSumario(w http.ResponseWriter, r *http.Request) {}

func (a *App) GetErrosQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) CreateQuestao(w http.ResponseWriter, r *http.Request) {

	user, err := auth.VerifyIdToken(r.Header.Get("Authorization"))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err = user.Get(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if user.NivelAcesso < 1 {
		utils.RespondWithError(w, http.StatusUnauthorized, "Usuário não autorizado a realizar a operação.")
		return
	}

	var q models.Questao
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

func (a *App) ReportQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) UpdateQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) DeleteQuestao(w http.ResponseWriter, r *http.Request) {}
