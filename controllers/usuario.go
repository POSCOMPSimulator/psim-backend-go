package controllers

import (
	"net/http"

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

	if err = user.Get(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// TODO: Adicionar estatísticas

	utils.RespondWithJSON(w, http.StatusOK, user)
	return
}

func (a *App) PromoteUsuario(w http.ResponseWriter, r *http.Request) {}

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
