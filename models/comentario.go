package models

import (
	"database/sql"
	"errors"
)

type Comentario struct {
	ID             int    `json:"id"`
	AutorID        int    `json:"autor_id"`
	AutorNome      string `json:"autor"`
	QuestaoID      int    `json:"questao_id,omitempty"`
	ProfilePicture string `json:"picture"`
	Texto          string `json:"texto"`
	DataPublicacao string `json:"data_publicacao"`
	Sinalizado     int    `json:"numero_sinalizacoes"`
}

type BatchComentarios struct {
	QuestaoID   int          `json:"-"`
	Comentarios []Comentario `json:"comentarios"`
}

func (bc *BatchComentarios) GetComentariosSinalizados(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (bc *BatchComentarios) GetComentariosQuestao(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (c *Comentario) Post(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (c *Comentario) Report(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (c *Comentario) Delete(db *sql.DB) error {
	return errors.New("Not implemented")
}
