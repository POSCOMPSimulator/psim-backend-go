package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) initializeRoutes() {

	public := a.Router.Group("/")

	authorized := a.Router.Group("/")
	authorized.Use(authMiddleware(a.tokenMaker, 0))

	moderator := a.Router.Group("/")
	moderator.Use(authMiddleware(a.tokenMaker, 1))

	// Rota base
	public.GET("/", func(ctx *gin.Context) { ctx.String(http.StatusOK, "PSIM Backend 2.0.0 in Golang") })

	// Rotas de usuário
	public.POST("/usuario/", a.CreateUsuario)
	public.POST("/usuario/login/", a.LoginUsuario)
	public.POST("/usuario/refresh", a.RenewTokenUsuario)
	authorized.GET("/usuario/", a.GetUsuario)
	authorized.DELETE("/usuario/", a.DeleteUsuario)
	moderator.PUT("/usuario/", a.PromoteUsuario)

	// // Rotas de questão
	public.GET("/questao/", a.GetQuestoes)
	public.GET("/questao/sumario", a.GetQSumario)
	public.PUT("/questao/", a.ReportQuestao)
	moderator.GET("/questao/:id/erros/", a.GetErrosQuestao)
	moderator.DELETE("/questao/:id/erros/", a.SolveErrosQuestao)
	moderator.POST("/questao/", a.CreateQuestao)
	moderator.PATCH("/questao/", a.UpdateQuestao)
	moderator.DELETE("/questao/:id", a.DeleteQuestao)

	// // Rotas de simulado
	authorized.GET("/simulado/", a.GetSimulados)
	authorized.POST("/simulado/", a.CreateSimulado)
	authorized.GET("/simulado/:id", a.GetSimulado)
	authorized.PUT("/simulado/:id/:to_state", a.UpdateStateSimulado)
	authorized.PATCH("/simulado/:id", a.UpdateRespostasSimulado)
	authorized.DELETE("/simulado/:id", a.DeleteSimulado)

	// // Rotas de comentários
	public.GET("/comentario/questao/:id", a.GetComentariosQuestao)
	public.PUT("/comentario/:id", a.ReportComentario)
	authorized.POST("/comentario/questao/:id", a.PostComentarioQuestao)
	authorized.DELETE("/comentario/:id", a.DeleteComentario)
	moderator.GET("/comentario/", a.GetComentariosSinalizados)
	moderator.DELETE("/comentario/:id/reports/", a.CleanComentario)

}
