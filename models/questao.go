package models

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
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
	Questoes []Questao     `json:"questoes"`
	Filtros  FiltroQuestao `json:"-"`
}

type FiltroQuestao struct {
	Anos        []int
	Areas       []string
	Sinalizadas bool
}

type ErrosQuestao struct {
	ID    int      `json:"-"`
	Erros []string `json:"erros"`
}

type MensagemErro struct {
	ID   int    `json:"-"`
	Erro string `json:"mensagem_erro"`
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

	bq.Questoes = []Questao{}
	ids := []interface{}{}
	map_id_questao := map[int]Questao{}

	queryString, args := bq.mountFilterQuery()
	rows, err := db.Query(queryString, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		q := Questao{}
		err := rows.Scan(&q.ID, &q.Ano, &q.Numero, &q.Area, &q.Subarea,
			&q.Alternativas[0], &q.Alternativas[1], &q.Alternativas[2],
			&q.Alternativas[3], &q.Alternativas[4], &q.Resposta, &q.Sinalizada)

		if err != nil {
			return err
		}

		q.ImagensQuestao.A = []string{}
		q.ImagensQuestao.B = []string{}
		q.ImagensQuestao.C = []string{}
		q.ImagensQuestao.D = []string{}
		q.ImagensQuestao.E = []string{}
		q.ImagensQuestao.Enunciado = []string{}

		ids = append(ids, q.ID)
		map_id_questao[q.ID] = q
	}

	if len(ids) > 0 {

		seps := map[bool]string{true: ",$", false: "$"}

		auxQuery := " WHERE id_questao IN ("
		for i := 0; i < len(ids); i++ {
			auxQuery += seps[i > 0] + strconv.Itoa(i+1)
		}
		auxQuery += ")"

		rows, err = db.Query("SELECT id_questao, texto FROM enunciado_questao"+auxQuery, ids...)
		if err != nil {
			return err
		}

		for rows.Next() {
			var (
				id    int
				texto string
			)

			rows.Scan(&id, &texto)

			q := map_id_questao[id]
			q.Enunciado = append(q.Enunciado, texto)
			map_id_questao[q.ID] = q

		}

		rows, err = db.Query("SELECT * FROM imagem_questao"+auxQuery, ids...)
		if err != nil {
			return err
		}

		for rows.Next() {
			var (
				id   int
				tipo string
				url  string
			)

			rows.Scan(&id, &tipo, &url)

			q := map_id_questao[id]

			switch tipo {
			case "A":
				q.ImagensQuestao.A = append(q.ImagensQuestao.A, url)
			case "B":
				q.ImagensQuestao.B = append(q.ImagensQuestao.B, url)
			case "C":
				q.ImagensQuestao.C = append(q.ImagensQuestao.C, url)
			case "D":
				q.ImagensQuestao.D = append(q.ImagensQuestao.D, url)
			case "E":
				q.ImagensQuestao.E = append(q.ImagensQuestao.E, url)
			default:
				q.ImagensQuestao.Enunciado = append(q.ImagensQuestao.Enunciado, url)
			}

			map_id_questao[q.ID] = q

		}

		for _, q := range map_id_questao {
			bq.Questoes = append(bq.Questoes, q)
		}

	}

	return nil
}

func (sq *SumarioQuestoes) Get(db *sql.DB) error {

	sq.Anos = []int{}
	sq.Areas = []string{}
	sq.Subareas = []string{}

	rows, err := db.Query("SELECT DISTINCT(ano) FROM questao")
	if err != nil {
		return errors.New("Não foi possível obter o sumário.")
	}

	for rows.Next() {
		var ano int
		rows.Scan(&ano)
		sq.Anos = append(sq.Anos, ano)
	}

	rows, err = db.Query("SELECT DISTINCT(area) FROM questao")
	if err != nil {
		return errors.New("Não foi possível obter o sumário.")
	}

	for rows.Next() {
		var area string
		rows.Scan(&area)
		sq.Areas = append(sq.Areas, area)
	}

	rows, err = db.Query("SELECT DISTINCT(subarea) FROM questao")
	if err != nil {
		return errors.New("Não foi possível obter o sumário.")
	}

	for rows.Next() {
		var subarea string
		rows.Scan(&subarea)
		sq.Subareas = append(sq.Subareas, subarea)
	}

	return nil
}

func (eq *ErrosQuestao) Get(db *sql.DB) error {

	eq.Erros = []string{}

	rows, err := db.Query("SELECT msg_err FROM sinalizacao_questao WHERE id_questao = $1", eq.ID)
	if err != nil {
		return errors.New("Não foi possível obter os erros.")
	}

	for rows.Next() {
		var msg_err string
		rows.Scan(&msg_err)
		eq.Erros = append(eq.Erros, msg_err)
	}

	return nil
}

func (eq *ErrosQuestao) Solve(db *sql.DB) error {
	return errors.New("Not implemented yet.")
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

func (m *MensagemErro) Report(db *sql.DB) error {

	if err := db.QueryRow("SELECT id FROM questao WHERE id = $1", m.ID).Scan(&m.ID); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Questão não foi encontrada.")
		}
		return err
	}

	if _, err := db.Exec("INSERT INTO sinalizacao_questao(id_questao, msg_err) VALUES($1, $2)", m.ID, m.Erro); err != nil {

		if err.Error() == `pq: duplicate key value violates unique constraint "sinalizacao_questao_pkey"` {
			return nil
		}

		return errors.New("Não foi possível reportar o erro.")
	}

	if _, err := db.Exec("UPDATE questao SET sinalizada = true WHERE id = $1", m.ID); err != nil {
		return errors.New("Não foi possível reportar o erro.")
	}

	return nil
}

func (q *Questao) Update(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (q *Questao) Delete(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (bq *BatchQuestoes) mountFilterQuery() (string, []interface{}) {

	seps := map[bool]string{true: ",$", false: "$"}

	max_len := len(bq.Filtros.Anos) + len(bq.Filtros.Areas)

	args, ind_args, ind_query := make([]interface{}, max_len), 0, 1
	queries := []string{"SELECT * FROM questao"}

	if len(bq.Filtros.Anos) > 0 {

		queryFilter := "ano IN ("

		for i, v := range bq.Filtros.Anos {
			args[ind_args] = v
			ind_args++
			queryFilter += seps[i > 0] + strconv.Itoa(ind_args)
		}

		queryFilter += ")"
		queries = append(queries, queryFilter)
		ind_query++

	}

	if len(bq.Filtros.Areas) > 0 {

		queryFilter := "area IN ("

		for i, v := range bq.Filtros.Areas {
			args[ind_args] = v
			ind_args++
			queryFilter += seps[i > 0] + strconv.Itoa(ind_args)
		}

		queryFilter += ")"
		queries = append(queries, queryFilter)
		ind_query++

	}

	if bq.Filtros.Sinalizadas {
		queries = append(queries, "sinalizada = true")
		ind_query++
	}

	if len(queries) > 1 {
		queries[0] += " WHERE "
		queries[0] += strings.Join(queries[1:], " AND ")
	}

	return queries[0] + " ORDER BY id", args

}
