package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) initializeRoutes() {

	public := a.Router.Group("/")

	// Rota base
	public.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "PSIM Backend 2.0.0 in Golang") })

	// // Rotas de questão
	public.GET("/questao/", a.GetQuestoes)
	public.GET("/questao/sumario/", a.GetQSumario)
	public.POST("/questao/:admincode/", a.CreateQuestao)

	// // Rotas de simulado
	public.POST("/simulado/", a.CreateSimulado)
	public.GET("/simulado/:id/", a.GetSimulado)
	public.PUT("/simulado/:id/:to_state/", a.UpdateStateSimulado)
	public.PATCH("/simulado/:id/", a.UpdateRespostasSimulado)

}
