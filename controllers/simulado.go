package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) CreateSimulado(ctx *gin.Context) {

	var err error
	var sim models.Simulado
	sim.ID, err = a.SimIDGen.Generate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("não foi possível criar o simulado")))
		return
	}

	sim.Estado = 0

	if err := ctx.ShouldBindJSON(&sim); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	if err := sim.Create(a.DB); err != nil {
		ctx.JSON(http.StatusNotAcceptable, utils.RespondWithError(err))
		return
	}

	ctx.JSON(http.StatusCreated, map[string]string{"text": sim.ID})

}

func (a *App) GetStateSimulado(ctx *gin.Context) {

	var sim models.Simulado

	sim.ID = ctx.Param("id")

	if err := sim.GetEstado(a.DB); err != nil {
		ctx.JSON(http.StatusNotFound, utils.RespondWithError(err))
		return
	}

	ctx.JSON(http.StatusOK, sim)

}

func (a *App) GetSimulado(ctx *gin.Context) {

	var sim models.Simulado

	sim.ID = ctx.Param("id")

	if err := sim.Get(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.JSON(http.StatusOK, sim)

}

func (a *App) UpdateStateSimulado(ctx *gin.Context) {

	var err error
	var sim models.Simulado

	sim.ID = ctx.Param("id")

	to_state := ctx.Param("to_state")
	switch strings.ToUpper(to_state) {
	case "INICIAR":

		if err := sim.Start(a.DB); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
			return
		}

		var retsim models.Simulado
		retsim.Questoes = sim.Questoes
		retsim.TempoRestante = sim.TempoRestante
		retsim.Respostas = sim.Respostas

		ctx.JSON(http.StatusAccepted, retsim)
		return

	case "CONTINUAR":

		if err := sim.Continue(a.DB); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
			return
		}

		var retsim models.Simulado
		retsim.Questoes = sim.Questoes
		retsim.TempoRestante = sim.TempoRestante
		retsim.Respostas = sim.Respostas

		ctx.JSON(http.StatusAccepted, retsim)
		return

	case "FINALIZAR":

		if err = sim.Finish(a.DB); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
			return
		}

		ctx.Status(http.StatusAccepted)
		return

	default:

		ctx.Status(http.StatusMethodNotAllowed)
		return

	}

}

func (a *App) UpdateRespostasSimulado(ctx *gin.Context) {

	var sim models.Simulado

	sim.ID = ctx.Param("id")

	var bresp models.BatchRespostas
	bresp.IDSimulado = sim.ID

	if err := ctx.ShouldBindJSON(&bresp); err != nil {

		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	if err := bresp.Update(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.Status(http.StatusOK)

}
