package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) CreateOrLoginUsuario(w http.ResponseWriter, r *http.Request) {

	authorization := strings.Fields(r.Header.Get("Authorization"))

	if len(authorization) != 2 {
		utils.RespondWithError(w, http.StatusAccepted, "erro")
	}

	payload, valid, _ := auth.VerifyIdToken(authorization[1])

	if valid {
		fmt.Println(payload)
	}

}

func (a *App) GetUsuario(w http.ResponseWriter, r *http.Request) {}

func (a *App) PromoteUsuario(w http.ResponseWriter, r *http.Request) {}

func (a *App) DeleteUsuario(w http.ResponseWriter, r *http.Request) {}
