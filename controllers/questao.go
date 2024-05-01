package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"poscomp-simulator.com/backend/models/questao"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetQuestoes(ctx *gin.Context) {

	type queryRequest struct {
		Anos        []int    `form:"anos"`
		Areas       []string `form:"areas"`
		Subareas    []string `form:"subareas"`
		Sinalizadas bool     `form:"sinalizadas"`
	}

	query := queryRequest{}
	if err := ctx.ShouldBind(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	var batch questao.BatchQuestoes
	batch.Filtros.Areas = query.Areas
	batch.Filtros.Anos = query.Anos
	batch.Filtros.Sinalizadas = query.Sinalizadas
	batch.Filtros.Subareas = query.Subareas

	if err := batch.Get(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.JSON(http.StatusOK, batch)

}

func (a *App) GetQSumario(ctx *gin.Context) {

	var sq questao.SumarioQuestoes
	sq.Get(a.DB)
	ctx.JSON(http.StatusOK, sq)

}

func (a *App) CreateQuestao(ctx *gin.Context) {

	admin_code := ctx.Param("admincode")

	if admin_code != a.AdminCode {
		err := errors.New("sem autorização")
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(err))
		return
	}

	var q questao.Questao
	if err := ctx.ShouldBindJSON(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	if err := q.Create(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.Status(http.StatusCreated)

}
