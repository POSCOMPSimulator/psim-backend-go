package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"poscomp-simulator.com/backend/models/questao"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetQuestoes(ctx *gin.Context) {

	type queryRequest struct {
		Anos        []int    `form:"anos"`
		Areas       []string `form:"areas"`
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

func (a *App) GetErrosQuestao(ctx *gin.Context) {

	var errosq questao.ErrosQuestao
	var err error
	qid := ctx.Param("id")

	errosq.ID, err = strconv.Atoi(qid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	errosq.Get(a.DB)
	ctx.JSON(http.StatusOK, errosq)

}

func (a *App) SolveErrosQuestao(ctx *gin.Context) {

	var errosq questao.ErrosQuestao
	var err error
	qid := ctx.Param("id")

	errosq.ID, err = strconv.Atoi(qid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	if err := ctx.ShouldBindJSON(&errosq); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	errosq.Solve(a.DB)

}

func (a *App) CreateQuestao(ctx *gin.Context) {

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

func (a *App) ReportQuestao(ctx *gin.Context) {

	var m questao.MensagemErro
	if err := ctx.ShouldBindJSON(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	if err := m.Report(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

}

func (a *App) UpdateQuestao(ctx *gin.Context) {

	var q questao.Questao
	if err := ctx.ShouldBindJSON(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	if err := q.Update(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

}

func (a *App) DeleteQuestao(ctx *gin.Context) {

	var err error
	var q questao.Questao
	qid := ctx.Param("id")
	q.ID, err = strconv.Atoi(qid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	q.Delete(a.DB)

}
