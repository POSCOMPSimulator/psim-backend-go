package models

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"poscomp-simulator.com/backend/models/questao"
)

const tempoMinimoPQuestao int = 3
const tempoMaximoPQuestao int = 5

type Simulado struct {
	questao.BatchQuestoes

	ID                int             `json:"id,omitempty"`
	Nome              string          `json:"nome,omitempty"`
	Estado            int             `json:"estado,omitempty"`
	TempoLimite       int             `json:"tempo_limite,omitempty"`
	TempoRestante     int             `json:"tempo_restante,omitempty"`
	IdUsuario         string          `json:"id_usuario,omitempty"`
	Anos              []int           `json:"anos,omitempty"`
	Areas             []string        `json:"areas,omitempty"`
	NumeroQuestoes    *NumeroQuestoes `json:"numero_questoes,omitempty"`
	Correcao          *Correcao       `json:"correcao,omitempty"`
	Respostas         Respostas       `json:"respostas_atuais,omitempty"`
	ContinuarSimulado bool            `json:"-"`
	Finalizar         bool            `json:"-"`
}

type BatchSimulados struct {
	IDUsuario int        `json:"-"`
	Simulados []Simulado `json:"simulados"`
}

type BatchRespostas struct {
	IDSimulado    int       `json:"-"`
	Respostas     Respostas `json:"respostas"`
	TempoRestante int       `json:"tempo_restante"`
}

type NumeroQuestoes struct {
	Tot int `json:"tot,omitempty"`
	Mat int `json:"mat"`
	Fun int `json:"fun"`
	Tec int `json:"tec"`
}

type Respostas map[int]string

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

func (s *Simulado) Start(db *sql.DB) error {

	if err := s.getEstado(db); err != nil {
		return err
	}

	if s.Estado != 0 {
		return errors.New("Simulado já foi iniciado.")
	}

	if _, err := db.Exec("UPDATE simulado SET estado = 1 WHERE id = $1", s.ID); err != nil {
		return errors.New("Não foi possível criar simulado")
	}

	if err := s.setQuestoes(db); err != nil {
		return err
	}

	s.getQuestoes(db)
	s.getRespostas(db)

	return nil
}

func (s *Simulado) Continue(db *sql.DB) error {

	if err := s.getEstado(db); err != nil {
		return err
	}

	if s.Estado != 1 {
		return errors.New("Simulado não foi iniciado ou está finalizado.")
	}

	s.getQuestoes(db)
	s.getRespostas(db)

	return nil
}

func (s *Simulado) Finish(db *sql.DB) error {

	if err := s.getEstado(db); err != nil {
		return err
	}

	if s.Estado != 1 {
		return errors.New("Simulado não foi iniciado ou está finalizado.")
	}

	if _, err := db.Exec("UPDATE simulado SET estado = 2 WHERE id = $1", s.ID); err != nil {
		return errors.New("Não foi possível criar simulado")
	}

	return nil

}

func (s *Simulado) Delete(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (s *Simulado) getEstado(db *sql.DB) error {

	s.NumeroQuestoes = &NumeroQuestoes{}
	var user string

	if err := db.QueryRow("SELECT id_usuario, estado, tempo_restante FROM simulado WHERE id = $1", s.ID).
		Scan(&user, &s.Estado, &s.TempoRestante); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Simulado não encontrado.")
		}
		return errors.New("Não foi recuperar o estado do simulado")
	}

	if user != s.IdUsuario {
		return errors.New("Simulado não pertence ao usuário.")
	}

	return nil

}

func (s *Simulado) setQuestoes(db *sql.DB) error {

	var (
		qmat int
		qfun int
		qtec int
	)

	if err := db.QueryRow("SELECT quant_mat, quant_fun, quant_tec FROM simulado WHERE id = $1", s.ID).Scan(&qmat, &qfun, &qtec); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Simulado não encontrado.")
		}
		return errors.New("Não foi recuperar o estado do simulado")
	}

	queryString := `
	SELECT id FROM(
	(SELECT id, numero, ano FROM questao
	WHERE area = 'Matemática'
	ORDER BY RANDOM()
	LIMIT $1)
	UNION
	(SELECT id, numero, ano FROM questao
	WHERE area = 'Fundamentos da Computação'
	ORDER BY RANDOM()
	LIMIT $2)
	UNION
	(SELECT id, numero, ano FROM questao
	WHERE area = 'Tecnologia da Computação'
	ORDER BY RANDOM()
	LIMIT $3)
	ORDER BY numero, ano) AS qsimulado
	`

	rows, err := db.Query(queryString, qmat, qfun, qtec)
	stmt, err2 := db.Prepare("INSERT INTO questoes_simulado(id_simulado, id_usuario, id_questao) VALUES ($1, $2, $3)")
	defer stmt.Close()

	if err != nil || err2 != nil {
		return errors.New("Não foi possível selecionar as questões.")
	}

	for rows.Next() {
		var idq int
		rows.Scan(&idq)
		stmt.Exec(s.ID, s.IdUsuario, idq)
	}

	return nil

}

func (s *Simulado) getQuestoes(db *sql.DB) error {

	s.Questoes = []questao.Questao{}

	query := `
	SELECT questao.*
	FROM questoes_simulado
	LEFT JOIN questao ON questao.id = questoes_simulado.id_questao
	WHERE id_simulado = $1`
	args := []interface{}{s.ID}

	s.SelectQuestoes(db, query, args)

	sort.Slice(s.Questoes, func(i, j int) bool {
		if s.Questoes[i].Numero < s.Questoes[j].Numero {
			return true
		} else if s.Questoes[i].Numero == s.Questoes[j].Numero {
			if s.Questoes[i].Ano < s.Questoes[j].Ano {
				return true
			}
		}
		return false
	})

	return nil

}

func (s *Simulado) getRespostas(db *sql.DB) error {

	s.Respostas = map[int]string{}
	rows, err := db.Query("SELECT id_questao, resposta FROM questoes_simulado WHERE id_simulado = $1", s.ID)
	if err != nil {
		return errors.New("Não foi possível obter as respostas.")
	}

	for rows.Next() {
		var id int
		var resp string
		rows.Scan(&id, &resp)
		s.Respostas[id] = resp
	}

	return nil
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
