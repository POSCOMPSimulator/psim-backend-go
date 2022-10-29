package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetComentariosSinalizados(ctx *gin.Context) {

	var bc models.BatchComentarios
	if err := bc.GetComentariosSinalizados(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.JSON(http.StatusOK, bc)

}

func (a *App) GetComentariosQuestao(ctx *gin.Context) {

	var err error
	var bc models.BatchComentarios

	id := ctx.Param("id")

	bc.QuestaoID, err = strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	if err := bc.GetComentariosQuestao(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.JSON(http.StatusOK, bc)

}

func (a *App) PostComentarioQuestao(ctx *gin.Context) {

	type postComment struct {
		Texto          string `json:"texto" binding:"required"`
		DataPublicacao string `json:"data_publicacao" binding:"required"`
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	qid := ctx.Param("id")

	req := postComment{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	c := models.Comentario{
		AutorID:        authPayload.UserID,
		Texto:          req.Texto,
		DataPublicacao: req.DataPublicacao,
	}

	var err error
	c.QuestaoID, err = strconv.Atoi(qid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	if err := c.Post(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.Status(http.StatusCreated)

}

func (a *App) ReportComentario(ctx *gin.Context) {

	var (
		c   models.Comentario
		err error
	)

	cid := ctx.Param("id")

	c.ID, err = strconv.Atoi(cid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	if err := c.Report(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.Status(http.StatusCreated)

}

func (a *App) CleanComentario(ctx *gin.Context) {

	var c models.Comentario
	cid := ctx.Param("id")

	var err error
	c.ID, err = strconv.Atoi(cid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	if err := c.Clean(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.Status(http.StatusOK)

}

func (a *App) DeleteComentario(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)
	var c models.Comentario
	cid := ctx.Param("id")

	var err error
	c.ID, err = strconv.Atoi(cid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	if err := c.Get(a.DB); err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.RespondWithError(errors.New("Comentário não foi encontrado.")))
			return
		}
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	if authPayload.UserID == c.AutorID || authPayload.UserLevel > 0 {
		c.Delete(a.DB)
		ctx.Status(http.StatusOK)
		return
	}

	ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(errors.New("Comentário não pertence ao usuário ou nível de acesso insuficiente.")))
	return

}
