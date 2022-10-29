package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"poscomp-simulator.com/backend/auth"
	"poscomp-simulator.com/backend/models"
	"poscomp-simulator.com/backend/utils"
)

func (a *App) GetSimulados(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)

	var bsim models.BatchSimulados
	bsim.IDUsuario = authPayload.UserID

	if err := bsim.Get(a.DB); err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.RespondWithError(err))
		return
	}

	ctx.JSON(http.StatusOK, bsim)

}

func (a *App) CreateSimulado(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)

	var sim models.Simulado
	sim.IdUsuario = authPayload.UserID
	sim.Estado = 0

	if err := ctx.ShouldBindJSON(&sim); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
		return
	}

	if err := sim.Create(a.DB); err != nil {
		ctx.JSON(http.StatusNotAcceptable, utils.RespondWithError(err))
		return
	}

	ctx.Status(http.StatusCreated)

}

func (a *App) GetSimulado(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)

	var err error
	var sim models.Simulado
	sim.IdUsuario = authPayload.UserID

	sim.ID, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	if err := sim.Get(a.DB); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
		return
	}

	ctx.JSON(http.StatusOK, sim)

}

func (a *App) UpdateStateSimulado(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)

	var err error
	var sim models.Simulado
	sim.IdUsuario = authPayload.UserID

	sim.ID, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

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

		var bresp models.BatchRespostas
		bresp.IDSimulado = sim.ID
		bresp.IDUsuario = authPayload.UserID

		if err := ctx.ShouldBindJSON(&bresp); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.RespondValidationError(err))
			return
		}

		if err := bresp.Update(a.DB); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
			return
		}

		if err = sim.Finish(a.DB); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.RespondWithError(err))
			return
		}

		ctx.Status(http.StatusMovedPermanently)

	default:

		ctx.Status(http.StatusMethodNotAllowed)
		return

	}

}

func (a *App) UpdateRespostasSimulado(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)

	var err error
	var sim models.Simulado
	sim.IdUsuario = authPayload.UserID

	sim.ID, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	var bresp models.BatchRespostas
	bresp.IDSimulado = sim.ID
	bresp.IDUsuario = authPayload.UserID

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

func (a *App) DeleteSimulado(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*auth.Payload)

	var err error
	var sim models.Simulado
	sim.IdUsuario = authPayload.UserID

	sim.ID, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.RespondWithError(errors.New("ID mal formatado.")))
		return
	}

	sim.Delete(a.DB)
	ctx.Status(http.StatusOK)

}
