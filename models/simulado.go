package models

import (
	"database/sql"
	"errors"
	"fmt"

	"poscomp-simulator.com/backend/models/questao"
)

const tempoMinimoPQuestao int = 3
const tempoMaximoPQuestao int = 5

type Simulado struct {
	questao.BatchQuestoes

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

	s.NumeroQuestoes.Tot = s.NumeroQuestoes.Mat + s.NumeroQuestoes.Fun + s.NumeroQuestoes.Tec

	err := db.QueryRow("SELECT id FROM simulado WHERE nome = $1", s.Nome).Scan(&s.ID)

	if err == nil {
		return errors.New("Simulado de mesmo nome já foi criado.")
	} else if err != sql.ErrNoRows {
		return errors.New("Não foi possível criar o simulado.")
	}

	if !(tempoMinimoPQuestao*s.NumeroQuestoes.Tot <= s.TempoLimite && s.TempoLimite <= tempoMaximoPQuestao*s.NumeroQuestoes.Tot) {
		return errors.New("Tempo limite fora do intervalo ideal.")
	}

	numeroMaximoQuestoes := getNumeroMaximoQuestoes(db, s.Anos)

	for _, area := range s.Areas {

		var num int

		switch area {
		case "Matemática":
			num = s.NumeroQuestoes.Mat
		case "Fundamentos da Computação":
			num = s.NumeroQuestoes.Fun
		case "Tecnologia da Computação":
			num = s.NumeroQuestoes.Tec
		}

		if num > numeroMaximoQuestoes[area] {
			return errors.New("Número de questões da área " + area + " ultrapassa o limite disponível.")
		}

	}

	query := `
	INSERT INTO 
	simulado(nome, estado, tempo_limite, 
			 quant_tot, quant_mat, quant_fun, 
			 quant_tec, tempo_restante, id_usuario)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if _, err := db.Exec(query, s.Nome, s.Estado, s.TempoLimite,
		s.NumeroQuestoes.Tot, s.NumeroQuestoes.Mat, s.NumeroQuestoes.Fun,
		s.NumeroQuestoes.Tec, s.TempoLimite, s.IdUsuario); err != nil {
		fmt.Println(err)
		return errors.New("Não foi possível criar o simulado.")
	}

	return nil

}

func (s *Simulado) Get(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (s *Simulado) Delete(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getNumeroMaximoQuestoes(db *sql.DB, anos []int) map[string]int {

	qtdQuestoes := map[string]int{}
	query := "SELECT ano, area, count(id) FROM questao GROUP BY ano, area"
	rows, _ := db.Query(query)

	for rows.Next() {
		var (
			ano  int
			area string
			qtd  int
		)

		rows.Scan(&ano, &area, &qtd)
		for _, v := range anos {
			if ano == v {
				qtdQuestoes[area] += qtd
				break
			}
		}
	}

	return qtdQuestoes
}
