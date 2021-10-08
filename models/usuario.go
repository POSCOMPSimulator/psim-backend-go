package models

import (
	"database/sql"
	"errors"
	"os"
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

func (u *Usuario) GetDummy() {
	u.Email = os.Getenv("DUMMY_TOKEN")
	u.FotoPerfil = "dummy.png"
	u.Nome = "Dummy User"
	u.NivelAcesso = 32767
}

func (u *Usuario) Create(db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO usuario(email, nome, foto_perfil, nivel_acesso) VALUES($1, $2, $3, $4)", u.Email, u.Nome, u.FotoPerfil, u.NivelAcesso); err != nil {
		return errors.New("Usuário não pode ser criado.")
	}
	return nil
}

func (u *Usuario) Get(db *sql.DB) error {
	if err := db.QueryRow("SELECT email, nome, foto_perfil, nivel_acesso FROM usuario WHERE email=$1", u.Email).Scan(&u.Email, &u.Nome, &u.FotoPerfil, &u.NivelAcesso); err != nil {
		return errors.New("Usuário não encontrado.")
	}
	return nil
}

func (u *Usuario) Promote(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (u *Usuario) Delete(db *sql.DB) error {
	return errors.New("Not implemented")
}
