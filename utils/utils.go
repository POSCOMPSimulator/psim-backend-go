package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"poscomp-simulator.com/backend/models"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithText(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"text": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func AuthUser(db *sql.DB, w http.ResponseWriter, r *http.Request, minLevel int16) (bool, models.Usuario) {

	// user, err := auth.VerifyIdToken(r.Header.Get("Authorization"))
	// if err != nil {
	// 	RespondWithError(w, http.StatusUnauthorized, err.Error())
	// 	return false, models.Usuario{}
	// }

	// if err = user.Get(db); err != nil {
	// 	RespondWithError(w, http.StatusNotFound, err.Error())
	// 	return false, models.Usuario{}
	// }

	// if user.NivelAcesso < minLevel {
	// 	RespondWithError(w, http.StatusUnauthorized, "Usuário não autorizado a realizar a operação.")
	// 	return false, models.Usuario{}
	// }

	// return true, user

	return false, models.Usuario{}
}
