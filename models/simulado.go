package models

import (
	"database/sql"
	"errors"
)

type Simulado struct {
	BatchQuestoes

	ID             int            `json:"id,omitempty"`
	Nome           string         `json:"nome"`
	Estado         int            `json:"estado"`
	TempoLimite    int            `json:"tempo_limite,omitempty"`
	TempoRestante  int            `json:"tempo_restante,omitempty"`
	IdUsuario      string         `json:"id_usuario,omitempty"`
	Anos           []int          `json:"anos"`
	Areas          []string       `json:"areas"`
	NumeroQuestoes NumeroQuestoes `json:"numero_questoes"`
	Correcao       Correcao       `json:"correcao,omitempty"`
	Respostas      []string       `json:"respostas_atuais,omitempty"`
}

type BatchSimulados struct {
	IDUsuario int        `json:"-"`
	Simulados []Simulado `json:"simulados"`
}

type BatchRespostas struct {
	IDSimulado    int        `json:"-"`
	Respostas     []Resposta `json:"respostas"`
	TempoRestante int        `json:"tempo_restante"`
}

type NumeroQuestoes struct {
	Tot int `json:"tot,omitempty"`
	Mat int `json:"mat"`
	Fun int `json:"fun"`
	Tec int `json:"tec"`
}

type Resposta struct {
	IDQuestao int    `json:"id_questao"`
	Resp      string `json:"resposta_questao"`
}

type Correcao struct {
	DataFinalizacao string         `json:"data_finalizacao"`
	TempoRealizacao int            `json:"tempo_realizacao"`
	Acertos         NumeroQuestoes `json:"acertos"`
	Erros           NumeroQuestoes `json:"erros"`
	Brancos         NumeroQuestoes `json:"brancos"`
}

func (bs *BatchSimulados) Get(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (br *BatchRespostas) UpdateRespostas(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (s *Simulado) Create(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (s *Simulado) Get(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (s *Simulado) Delete(db *sql.DB) error {
	return errors.New("Not implemented")
}
