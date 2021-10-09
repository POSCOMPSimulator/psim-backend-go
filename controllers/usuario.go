package controllers

import (
	"net/http"
	"strconv"

	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) CreateOrLoginUsuario(w http.ResponseWriter, r *http.Request) {

	user, err := auth.VerifyIdToken(r.Header.Get("Authorization"))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err = user.Get(a.DB); err != nil {
		user.Create(a.DB)
		utils.RespondWithJSON(w, http.StatusCreated, user)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
	return

}

func (a *App) GetUsuario(w http.ResponseWriter, r *http.Request) {

	user, err := auth.VerifyIdToken(r.Header.Get("Authorization"))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user.Completo = true
	if err = user.Get(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
	return
}

func (a *App) PromoteUsuario(w http.ResponseWriter, r *http.Request) {

	user, err := auth.VerifyIdToken(r.Header.Get("Authorization"))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err = user.Get(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	email_target := r.FormValue("email")
	if email_target == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Campo email não encontrado.")
		return
	}

	next_level_st := r.FormValue("nivel")
	if next_level_st == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Campo nivel não encontrado.")
		return
	}

	next_level, err := strconv.ParseInt(next_level_st, 10, 16)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Campo nivel mal formatado.")
		return
	}

	if user.NivelAcesso < 1 || user.NivelAcesso < int16(next_level) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Usuário não autorizado a realizar a promoção de conta.")
		return
	}

	userToPromote := models.Usuario{Email: email_target, NivelAcesso: int16(next_level)}
	if err = userToPromote.Promote(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Usuário a ser promovido não encontrado.")
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (a *App) DeleteUsuario(w http.ResponseWriter, r *http.Request) {

	user, err := auth.VerifyIdToken(r.Header.Get("Authorization"))
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err = user.Get(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	email_target := r.FormValue("email")
	if email_target == "" {
		user.Delete(a.DB)
		return
	}

	userToDelete := models.Usuario{Email: email_target}
	if err = userToDelete.Get(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Usuário a ser deletado não encontrado.")
		return
	}

	if user.NivelAcesso < userToDelete.NivelAcesso {
		utils.RespondWithError(w, http.StatusUnauthorized, "Usuário não autorizado a realizar a exclusão de conta.")
		return
	}

	userToDelete.Delete(a.DB)
	return

}
