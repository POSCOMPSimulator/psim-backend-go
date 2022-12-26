package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/utils"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationType       = "bearer"
	authorizationPayloadKey = "auth_payload"
)

func authMiddleware(tokenMaker auth.Maker, minLevel int16) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("Authorization header não encontrado.")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespondWithError(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := errors.New("Authorization header mal formatado.")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespondWithError(err))
			return
		}

		if strings.ToLower(fields[0]) != authorizationType {
			err := errors.New("Tipo de autenticação não suportado.")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespondWithError(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespondWithError(err))
			return
		}

		if payload.UserLevel < minLevel {
			err := errors.New("Usuário sem nível de acesso para a operação.")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.RespondWithError(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
