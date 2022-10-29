package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "email":
		return "This field should be formated as a valid email"
	case "alphanum":
		return "This field should be composed of alphanumeric characters"
	case "len":
		return "The lenght of this field should be equal to " + fe.Param()
	case "eqfield":
		return "This field should be equal to " + fe.Param()
	}
	return "Unknown error"
}

func RespondValidationError(err error) gin.H {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
		}
		return gin.H{"errors": out}
	}
	return gin.H{"errors": err}
}

func RespondWithError(err error) gin.H {
	return gin.H{"error": err.Error()}
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
