package models

import (
	"database/sql"
	"errors"
)

type Usuario struct {
	Email        string           `json:"email,omitempty"`
	NivelAcesso  int16            `json:"nivel_acesso"`
	Nome         string           `json:"nome"`
	FotoPerfil   string           `json:"foto_perfil"`
	Estatisticas EstaticasUsuario `json:"estatisticas,omitempty"`
}

type EstaticasUsuario struct {
	NumSimuladoFinalizado     int                       `json:"num_simulados_finalizados"`
	NumComentariosPublicados  int                       `json:"num_comentarios_publicados"`
	PorcentagemQuestoesFeitas PorcentagemQuestoesFeitas `json:"porcentagem_questoes_feitas"`
}

type PorcentagemQuestoesFeitas struct {
	Geral int `json:"geral"`
	Mat   int `json:"mat"`
	Fun   int `json:"fun"`
	Tec   int `json:"tec"`
}

func (u *Usuario) Create(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (u *Usuario) Get(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (u *Usuario) Promote(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (u *Usuario) Delete(db *sql.DB) error {
	return errors.New("Not implemented")
}
