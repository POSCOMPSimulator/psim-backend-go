package questao

import (
	"database/sql"
	"sort"
	"strconv"
	"strings"
)

type BatchQuestoes struct {
	Questoes []Questao       `json:"questoes,omitempty"`
	Filtros  SumarioQuestoes `json:"-"`
}

func (bq *BatchQuestoes) Get(db *sql.DB) error {

	bq.Questoes = []Questao{}

	queryString, args := bq.mountFilterQuery()
	bq.SelectQuestoes(db, queryString, args)

	sort.Slice(bq.Questoes, func(i, j int) bool {
		if bq.Questoes[i].Ano < bq.Questoes[j].Ano {
			return true
		} else if bq.Questoes[i].Ano == bq.Questoes[j].Ano {
			if bq.Questoes[i].Numero < bq.Questoes[j].Numero {
				return true
			}
		}
		return false
	})

	return nil
}

func (bq *BatchQuestoes) SelectQuestoes(db *sql.DB, query string, args []interface{}) error {

	ids := []interface{}{}
	map_id_questao := map[int]Questao{}

	rows, err := db.Query(query, args...)
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
