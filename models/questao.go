package models

import (
	"database/sql"
	"errors"
)

type Questao struct {
	ID                  int         `json:"id"`
	Numero              int         `json:"numero"`
	Ano                 int         `json:"ano"`
	Area                string      `json:"area"`
	Subarea             string      `json:"subarea"`
	Alternativas        [5]string   `json:"alternativas"`
	Resposta            int         `json:"resposta"`
	Enunciado           []string    `json:"enunciado"`
	ImagensEnunciado    []string    `json:"imagens_enunciado"`
	ImagensAlternativas [][5]string `json:"imagens_alternativas"`
	Sinalizada          bool        `json:"sinalizada"`
}

func (q *Questao) GetQuestoes(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) GetQSumario(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) GetQuestao(db *sql.DB) error {
	err := db.QueryRow("SELECT * FROM questao WHERE id=$1", q.ID).Scan(&q.ID, //WHERE id=$1", q.ID
		&q.Numero, &q.Ano, &q.Area, &q.Subarea, &q.Alternativas[0],
		&q.Alternativas[1], &q.Alternativas[2], &q.Alternativas[3],
		&q.Alternativas[4], &q.Resposta)

	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT texto FROM enunciado_questao WHERE id=$1 ORDER BY ordem", q.ID)

	if err != nil {
		return err
	}

	defer rows.Close()

	q.Enunciado = []string{}

	for rows.Next() {
		var texto string

		if err := rows.Scan(&texto); err != nil {
			return err
		}

		q.Enunciado = append(q.Enunciado, texto)
	}

	return nil
}

func (q *Questao) CreateQuestao(db *sql.DB) error {

	queryString := "INSERT INTO questao(numero, ano, area, subarea, alternativa_a, alternativa_b, alternativa_c, alternativa_d, alternativa_e, resposta) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"

	err := db.QueryRow(queryString,
		q.Numero, q.Ano, q.Area, q.Subarea, q.Alternativas[0],
		q.Alternativas[1], q.Alternativas[2], q.Alternativas[3],
		q.Alternativas[4], q.Resposta).Scan(&q.ID)

	if err != nil {
		return err
	}

	return nil
}

func (q *Questao) ReportQuestao(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) UpdateQuestao(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) DeleteQuestao(db *sql.DB) error {
	return errors.New("Not implemented")
}
