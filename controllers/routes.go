package controllers

import (
	"net/http"

	"poscomp-simulator.com/backend/utils"
)

func (a *App) initializeRoutes() {

	// Rota base
	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.RespondWithText(w, http.StatusOK, "PSIM Backend 2.0.0 in Golang")
	})

	// Rotas de usuário
	a.Router.HandleFunc("/usuario/", a.GetUsuario).Methods("GET")
	a.Router.HandleFunc("/usuario/", a.CreateOrLoginUsuario).Methods("POST")
	a.Router.HandleFunc("/usuario/", a.PromoteUsuario).Methods("PUT")
	a.Router.HandleFunc("/usuario/", a.DeleteUsuario).Methods("DELETE")

	// Rotas de questão
	a.Router.HandleFunc("/questao/", a.GetQuestoes).Methods("GET")
	a.Router.HandleFunc("/questao/", a.CreateQuestao).Methods("POST")
	a.Router.HandleFunc("/questao/", a.ReportQuestao).Methods("PUT")
	a.Router.HandleFunc("/questao/", a.UpdateQuestao).Methods("PATCH")
	a.Router.HandleFunc("/questao/sumario/", a.GetQSumario).Methods("GET")
	a.Router.HandleFunc("/questao/{id}/", a.DeleteQuestao).Methods("DELETE")
	a.Router.HandleFunc("/questao/{id}/erros/", a.GetErrosQuestao).Methods("GET")
	a.Router.HandleFunc("/questao/{id}/erros/", a.SolveErrosQuestao).Methods("DELETE")

	// Rotas de simulado
	a.Router.HandleFunc("/simulado/", a.GetSimulados).Methods("GET")
	a.Router.HandleFunc("/simulado/", a.CreateSimulado).Methods("POST")
	a.Router.HandleFunc("/simulado/{id}/", a.GetSimulado).Methods("GET")
	a.Router.HandleFunc("/simulado/{id}/", a.UpdateStateSimulado).Methods("PUT")
	a.Router.HandleFunc("/simulado/{id}/", a.UpdateRespostasSimulado).Methods("PATCH")
	a.Router.HandleFunc("/simulado/{id}/", a.DeleteSimulado).Methods("DELETE")

	// Rotas de comentários
	a.Router.HandleFunc("/comentario/", a.GetComentariosSinalizados).Methods("GET")
	a.Router.HandleFunc("/comentario/questao/{id}/", a.GetComentariosQuestao).Methods("GET")
	a.Router.HandleFunc("/comentario/questao/{id}/", a.PostComentarioQuestao).Methods("POST")
	a.Router.HandleFunc("/comentario/{id}/", a.ReportComentario).Methods("PUT")
	a.Router.HandleFunc("/comentario/{id}/", a.DeleteComentario).Methods("DELETE")

}
