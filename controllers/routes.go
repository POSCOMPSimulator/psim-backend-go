package controllers

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/usuario/", a.CreateOrLoginUsuario).Methods("POST")

	a.Router.HandleFunc("/questao/", a.GetQuestoes).Methods("GET")
	a.Router.HandleFunc("/questao/", a.CreateQuestao).Methods("POST")
	a.Router.HandleFunc("/questao/sumario/", a.GetQSumario).Methods("GET")
	a.Router.HandleFunc("/questao/{id}/", a.GetQuestao).Methods("GET")
	a.Router.HandleFunc("/questao/{id}/", a.ReportQuestao).Methods("PUT")
	a.Router.HandleFunc("/questao/{id}/", a.UpdateQuestao).Methods("PATCH")
	a.Router.HandleFunc("/questao/{id}/", a.DeleteQuestao).Methods("DELETE")

}
