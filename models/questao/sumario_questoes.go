package questao

import (
	"database/sql"
	"errors"
)

type SumarioQuestoes struct {
	Anos        []int    `json:"anos"`
	Areas       []string `json:"areas"`
	Subareas    []string `json:"subareas"`
	Sinalizadas bool     `json:"-"`
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
