package models

import (
	"database/sql"
	"errors"
)

type Questao struct {
	ID             int            `json:"id"`
	Numero         int            `json:"numero"`
	Ano            int            `json:"ano"`
	Area           string         `json:"area"`
	Subarea        string         `json:"subarea"`
	Alternativas   [5]string      `json:"alternativas"`
	Resposta       int            `json:"resposta"`
	Enunciado      []string       `json:"enunciado"`
	ImagensQuestao ImagensQuestao `json:"imagens"`
	Sinalizada     bool           `json:"sinalizada"`
}

type BatchQuestoes struct {
	Questoes []Questao              `json:"questoes"`
	Filtros  map[string]interface{} `json:"-"`
}

type ErrosQuestao struct {
	Erros []string `json:"erros"`
}

type ImagensQuestao struct {
	Enunciado []string `json:"enunciado"`
	A         []string `json:"alternativa_a"`
	B         []string `json:"alternativa_b"`
	C         []string `json:"alternativa_c"`
	D         []string `json:"alternativa_d"`
	E         []string `json:"alternativa_e"`
}

type SumarioQuestoes struct {
	Anos     []int    `json:"anos"`
	Areas    []string `json:"areas"`
	Subareas []string `json:"subareas"`
}

func (bq *BatchQuestoes) Get(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (sq *SumarioQuestoes) Get(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) GetErros(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) Create(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) Report(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) Update(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) Delete(db *sql.DB) error {
	return errors.New("Not implemented")
}
