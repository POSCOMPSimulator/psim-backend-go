package questao

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

type ImagensQuestao struct {
	Enunciado []string `json:"enunciado"`
	A         []string `json:"alternativa_a"`
	B         []string `json:"alternativa_b"`
	C         []string `json:"alternativa_c"`
	D         []string `json:"alternativa_d"`
	E         []string `json:"alternativa_e"`
}

func (q *Questao) Create(db *sql.DB) error {

	if err := db.QueryRow("SELECT id FROM questao WHERE ano = $1 AND numero = $2", q.Ano, q.Numero).Scan(&q.ID); err == nil {
		return errors.New("Questão já foi adicionada.")
	}

	var queries = [3]string{
		`INSERT INTO questao(ano, numero, area, subarea, alternativa_a, alternativa_b, alternativa_c, alternativa_d, alternativa_e, gabarito)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
		`,
		"INSERT INTO enunciado_questao(id_questao, ordem, texto) VALUES($1, $2, $3)",
		"INSERT INTO imagem_questao(id_questao, tipo, url_img) VALUES($1, $2, $3)",
	}

	if err := db.QueryRow(queries[0], q.Ano, q.Numero, q.Area,
		q.Subarea, q.Alternativas[0], q.Alternativas[1],
		q.Alternativas[2], q.Alternativas[3],
		q.Alternativas[4], q.Resposta).Scan(&q.ID); err != nil {
		return errors.New("Questão não pode ser criada.")
	}

	for e, v := range q.Enunciado {
		if _, err := db.Exec(queries[1], q.ID, e, v); err != nil {
			return errors.New("Questão não pode ser criada.")
		}
	}

	for _, v := range q.ImagensQuestao.Enunciado {
		if _, err := db.Exec(queries[2], q.ID, "", v); err != nil {
			return errors.New("Questão não pode ser criada.")
		}
	}

	for _, v := range q.ImagensQuestao.A {
		if _, err := db.Exec(queries[2], q.ID, "A", v); err != nil {
			return errors.New("Questão não pode ser criada.")
		}
	}

	for _, v := range q.ImagensQuestao.B {
		if _, err := db.Exec(queries[2], q.ID, "B", v); err != nil {
			return errors.New("Questão não pode ser criada.")
		}
	}

	for _, v := range q.ImagensQuestao.C {
		if _, err := db.Exec(queries[2], q.ID, "C", v); err != nil {
			return errors.New("Questão não pode ser criada.")
		}
	}

	for _, v := range q.ImagensQuestao.D {
		if _, err := db.Exec(queries[2], q.ID, "D", v); err != nil {
			return errors.New("Questão não pode ser criada.")
		}
	}

	for _, v := range q.ImagensQuestao.E {
		if _, err := db.Exec(queries[2], q.ID, "E", v); err != nil {
			return errors.New("Questão não pode ser criada.")
		}
	}

	return nil
}

func (q *Questao) Update(db *sql.DB) error {

	var queries = [3]string{
		`UPDATE questao
		 SET subarea = $1, alternativa_a = $2, alternativa_b = $3, alternativa_c = $4, alternativa_d = $5, alternativa_e = $6, gabarito = $7
		 WHERE id = $8
		`,
		"INSERT INTO enunciado_questao(id_questao, ordem, texto) VALUES($1, $2, $3)",
		"INSERT INTO imagem_questao(id_questao, tipo, url_img) VALUES($1, $2, $3)",
	}

	if _, err := db.Exec(queries[0], q.Subarea, q.Alternativas[0], q.Alternativas[1],
		q.Alternativas[2], q.Alternativas[3], q.Alternativas[4], q.Resposta, q.ID); err != nil {
		return errors.New("Questão não pôde ser editada.")
	}

	if _, err := db.Exec("DELETE FROM enunciado_questao WHERE id_questao = $1", q.ID); err != nil {
		return errors.New("Questão não pôde ser editada.")
	}

	if _, err := db.Exec("DELETE FROM imagem_questao WHERE id_questao = $1", q.ID); err != nil {
		return errors.New("Questão não pôde ser editada.")
	}

	for e, v := range q.Enunciado {
		if _, err := db.Exec(queries[1], q.ID, e, v); err != nil {
			return errors.New("Questão não pôde ser editada.")
		}
	}

	for _, v := range q.ImagensQuestao.Enunciado {
		if _, err := db.Exec(queries[2], q.ID, "", v); err != nil {
			return errors.New("Questão não pôde ser editada.")
		}
	}

	for _, v := range q.ImagensQuestao.A {
		if _, err := db.Exec(queries[2], q.ID, "A", v); err != nil {
			return errors.New("Questão não pôde ser editada.")
		}
	}

	for _, v := range q.ImagensQuestao.B {
		if _, err := db.Exec(queries[2], q.ID, "B", v); err != nil {
			return errors.New("Questão não pôde ser editada.")
		}
	}

	for _, v := range q.ImagensQuestao.C {
		if _, err := db.Exec(queries[2], q.ID, "C", v); err != nil {
			return errors.New("Questão não pôde ser editada.")
		}
	}

	for _, v := range q.ImagensQuestao.D {
		if _, err := db.Exec(queries[2], q.ID, "D", v); err != nil {
			return errors.New("Questão não pôde ser editada.")
		}
	}

	for _, v := range q.ImagensQuestao.E {
		if _, err := db.Exec(queries[2], q.ID, "E", v); err != nil {
			return errors.New("Questão não pôde ser editada.")
		}
	}

	return nil
}

func (q *Questao) Delete(db *sql.DB) error {

	if _, err := db.Exec("DELETE FROM questao WHERE id = $1", q.ID); err != nil {
		return errors.New("Não foi possível remover a questão.")
	}

	return nil

}
