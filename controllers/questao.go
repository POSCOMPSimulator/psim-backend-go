package controllers

import (
	"net/http"
)

func (a *App) GetQuestoes(w http.ResponseWriter, r *http.Request) {}

func (a *App) GetQSumario(w http.ResponseWriter, r *http.Request) {}

func (a *App) GetErrosQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) CreateQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) ReportQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) UpdateQuestao(w http.ResponseWriter, r *http.Request) {}

func (a *App) DeleteQuestao(w http.ResponseWriter, r *http.Request) {}
