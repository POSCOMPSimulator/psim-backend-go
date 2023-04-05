package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	verificationCode := utils.GenerateVerificationCode()
	recoverCode := utils.GenerateVerificationCode()

	user := &models.Usuario{
		Email:             req.Email,
		Senha:             req.Password,
		Nome:              req.Username,
		CodigoVerificacao: verificationCode,
		CodigoRecuperacao: recoverCode,
	}

	if err := user.Get(a.DB); err != nil {
		user.Create(a.DB)
		a.Mailer.SendVerificationMail([]string{user.Email}, verificationCode)
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

	type loginUserResponse struct {
		User                  *models.Usuario `json:"user"`
		SessionID             uuid.UUID       `json:"session_id"`
		AccessToken           string          `json:"access_token"`
		AccessTokenExpiresAt  time.Time       `json:"access_token_expires_at"`
		RefreshToken          string          `json:"refresh_token"`
		RefreshTokenExpiresAt time.Time       `json:"refresh_token_expires_at"`
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
	access_token, access_payload, _ := a.tokenMaker.CreateToken(user.Email, user.NivelAcesso, user.Verificado, tokenDuration)

	refreshTokenDuration, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	refresh_token, refresh_payload, _ := a.tokenMaker.CreateToken(user.Email, user.NivelAcesso, user.Verificado, refreshTokenDuration)

	session := &models.Session{
		ID:           refresh_payload.ID,
		Username:     user.Email,
		RefreshToken: refresh_token,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refresh_payload.ExperiesAt,
	}

	if err := session.CreateSession(a.DB); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.RespondWithError(err))
		return
	}

	res := loginUserResponse{
		User:                  user,
		SessionID:             session.ID,
		AccessToken:           access_token,
		AccessTokenExpiresAt:  access_payload.ExperiesAt,
		RefreshToken:          refresh_token,
		RefreshTokenExpiresAt: refresh_payload.ExperiesAt,
	}

	ctx.JSON(http.StatusOK, res)
	return

}

func (a *App) RenewTokenUsuario(ctx *gin.Context) {

	type renewTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	type renewTokenResponse struct {
		AccessToken          string    `json:"access_token"`
		AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	}

	req := renewTokenRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	refresh_payload, err := a.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(err))
		return
	}

	session := &models.Session{
		ID: refresh_payload.ID,
	}

	if err := session.GetSession(a.DB); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.RespondWithError(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.RespondWithError(err))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(errors.New("blocked session")))
		return
	}

	if session.Username != refresh_payload.UserID {
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(errors.New("incorrect session user")))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(errors.New("mismatched session token")))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(errors.New("expired session")))
		return
	}

	tokenDuration, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	access_token, access_payload, _ := a.tokenMaker.CreateToken(refresh_payload.UserID, refresh_payload.UserLevel, refresh_payload.Verificado, tokenDuration)

	res := renewTokenResponse{
		AccessToken:          access_token,
		AccessTokenExpiresAt: access_payload.ExperiesAt,
	}

	ctx.JSON(http.StatusOK, res)
	return

}

func (a *App) VerificaUsuario(ctx *gin.Context) {

	type verificaUsuarioRequest struct {
		CodigoVerificacao string `json:"codigo_verificacao" binding:"required"`
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)

	req := verificaUsuarioRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	userToVerify := models.Usuario{Email: authPayload.UserID, CodigoVerificacao: req.CodigoVerificacao}
	if err := userToVerify.Verify(a.DB); err != nil {
		ctx.JSON(http.StatusNotAcceptable, utils.RespondWithError(errors.New("Código de verificação inválido.")))
		return
	}

	ctx.Status(http.StatusOK)
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
