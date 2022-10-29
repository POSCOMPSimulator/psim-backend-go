package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) CreateUsuario(ctx *gin.Context) {

	type createUserRequest struct {
		Email           string `json:"email" binding:"required,email"`
		Username        string `json:"nome" binding:"required,alphanum"`
		Password        string `json:"senha" binding:"required,len=8"`
		ConfirmPassword string `json:"confirma_senha" binding:"required,eqfield=Password"`
	}

	req := createUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	user := &models.Usuario{
		Email: req.Email,
		Senha: req.Password,
		Nome:  req.Username,
	}

	if err := user.Get(a.DB); err != nil {
		user.Create(a.DB)
		ctx.Status(http.StatusCreated)
		return
	}

	ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("Usuário já existe")))
	return

}

func (a *App) LoginUsuario(ctx *gin.Context) {

	type loginUserRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"senha" binding:"required,len=8"`
	}

	req := loginUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	user := &models.Usuario{
		Email: req.Email,
	}

	if err := user.Get(a.DB); err != nil {
		ctx.JSON(http.StatusNotFound, utils.RespondWithError(errors.New("Usuário não encontrado.")))
		return
	}

	if err := auth.CheckPassword(req.Password, user.Senha); err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(errors.New("Senha incorreta.")))
		return
	}

	tokenDuration, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	access_token, _ := a.tokenMaker.CreateToken(user.Email, user.NivelAcesso, tokenDuration)
	user.TokenAcesso = access_token

	ctx.JSON(http.StatusOK, user)
	return

}

func (a *App) GetUsuario(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	user := models.Usuario{
		Email:    authPayload.UserID,
		Completo: true,
	}

	if err := user.Get(a.DB); err != nil {
		ctx.JSON(http.StatusNotFound, utils.RespondWithError(errors.New("Usuário não encontrado.")))
		return
	}

	ctx.JSON(http.StatusOK, user)
	return
}

func (a *App) PromoteUsuario(ctx *gin.Context) {

	type promoteUserRequest struct {
		Email     string `json:"email" binding:"required,email"`
		NextLevel int16  `json:"nivel" binding:"required"`
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)

	req := promoteUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	if authPayload.UserLevel <= req.NextLevel {
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(errors.New("Usuário não autorizado a realizar a promoção de conta.")))
		return
	}

	userToPromote := models.Usuario{Email: req.Email, NivelAcesso: req.NextLevel}
	if err := userToPromote.Promote(a.DB); err != nil {
		ctx.JSON(http.StatusNotFound, utils.RespondWithError(errors.New("Usuário a ser promovido não encontrado.")))
		return
	}

	ctx.Status(http.StatusOK)

}

func (a *App) DeleteUsuario(ctx *gin.Context) {

	type deleteUserRequest struct {
		Email string `json:"email"`
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)

	req := deleteUserRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	if req.Email == "" {
		user := models.Usuario{Email: authPayload.UserID}
		user.Delete(a.DB)
		ctx.Status(http.StatusOK)
		return
	}

	userToDelete := models.Usuario{Email: req.Email}
	if err := userToDelete.Get(a.DB); err != nil {
		ctx.JSON(http.StatusNotFound, utils.RespondWithError(errors.New("Usuário a ser apagado não encontrado.")))
		return
	}

	if authPayload.UserLevel < userToDelete.NivelAcesso {
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(errors.New("Usuário não autorizado a realizar a exclusão de conta.")))
		return
	}

	userToDelete.Delete(a.DB)
	ctx.Status(http.StatusOK)
	return

}
