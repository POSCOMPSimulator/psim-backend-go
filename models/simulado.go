package models

import (
	"database/sql"
	"errors"
)

type Simulado struct {
	ID             int            `json:"id"`
	Nome           string         `json:"nome"`
	Estado         int            `json:"estado"`
	TempoLimite    int            `json:"tempo_limite"`
	TempoRestante  int            `json:"tempo_restante"`
	IdUsuario      string         `json:"id_usuario"`
	Anos           []int          `json:"anos"`
	Areas          []string       `json:"areas"`
	NumeroQuestoes NumeroQuestoes `json:"numero_questoes"`
}

type BatchSimulados struct {
	IDUsuario int        `json:"-"`
	Simulados []Simulado `json:"simulados"`
}

type SimuladoRealizacao struct {
	BatchQuestoes

	ID            int      `json:"id"`
	Respostas     []string `json:"respostas_atuais"`
	TempoRestante int      `json:"tempo_restante"`
}

type SimuladoRespondido struct {
	SimuladoRealizacao
	Correcao Correcao `json:"correcao"`
}

type NumeroQuestoes struct {
	Tot int `json:"tot"`
	Mat int `json:"mat"`
	Fun int `json:"fun"`
	Tec int `json:"tec"`
}

type BatchRespostas struct {
	Respostas     []Resposta `json:"respostas"`
	TempoRestante int        `json:"tempo_restante"`
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

func (s *Simulado) Create(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (s *SimuladoRespondido) Get(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (s *Simulado) Delete(db *sql.DB) error {
	return errors.New("Not implemented")
}
