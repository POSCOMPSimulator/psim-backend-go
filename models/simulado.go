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
	IDUsuario string     `json:"-"`
	Simulados []Simulado `json:"simulados"`
}

type BatchRespostas struct {
	IDSimulado    int       `json:"-"`
	IDUsuario     string    `json:"-"`
	Respostas     Respostas `json:"respostas"`
	TempoRestante int       `json:"tempo_restante"`
}

type NumeroQuestoes struct {
	Tot int `json:"tot,omitempty"`
	Mat int `json:"mat"`
	Fun int `json:"fun"`
	Tec int `json:"tec"`
}

type Respostas map[int]int

type Correcao struct {
	ID              int            `json:"-"`
	DataFinalizacao string         `json:"data_finalizacao"`
	TempoRealizacao int            `json:"tempo_realizacao"`
	Acertos         NumeroQuestoes `json:"acertos"`
	Erros           NumeroQuestoes `json:"erros"`
	Brancos         NumeroQuestoes `json:"brancos"`
}

func (bs *BatchSimulados) Get(db *sql.DB) error {

	rows, err := db.Query("SELECT * FROM simulado WHERE id_usuario = $1", bs.IDUsuario)
	if err != nil {
		return errors.New("Não foi possível recuperar os simulados.")
	}

	bs.Simulados = []Simulado{}
	for rows.Next() {
		var sim Simulado
		sim.NumeroQuestoes = &NumeroQuestoes{}

		rows.Scan(&sim.ID, &sim.Nome, &sim.Estado, &sim.TempoLimite,
			&sim.NumeroQuestoes.Tot, &sim.NumeroQuestoes.Mat, &sim.NumeroQuestoes.Fun,
			&sim.NumeroQuestoes.Tec, &sim.TempoRestante, &sim.IdUsuario)

		if sim.Estado == 2 {

			if err := sim.getCorrecao(db); err != nil {
				return err
			}

		}

		bs.Simulados = append(bs.Simulados, sim)

	}

	return nil

}

func (br *BatchRespostas) Update(db *sql.DB) error {

	stmt, err := db.Prepare("UPDATE questoes_simulado SET resposta = $1 WHERE id_simulado = $2 AND id_usuario = $3 AND id_questao = $4")
	if err != nil {
		return errors.New("Não foi possível atualizar as respostas.")
	}

	for id, v := range br.Respostas {
		_, err = stmt.Exec(v, br.IDSimulado, br.IDUsuario, id)
		if err != nil {
			return errors.New("Não foi possível atualizar as respostas.")
		}
	}

	if _, err := db.Exec("UPDATE simulado SET tempo_restante = $1 WHERE id = $2 AND id_usuario = $3",
		br.TempoRestante, br.IDSimulado, br.IDUsuario); err != nil {
		return errors.New("Não foi possível atualizar as respostas.")
	}

	return nil
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

	if err := s.getEstado(db); err != nil {
		return err
	}

	if s.Estado != 2 {
		return errors.New("Simulado não foi finalizado.")
	}

	if err := db.QueryRow("SELECT * FROM simulado WHERE id = $1", s.ID).Scan(&s.ID, &s.Nome, &s.Estado, &s.TempoLimite,
		&s.NumeroQuestoes.Tot, &s.NumeroQuestoes.Mat, &s.NumeroQuestoes.Fun,
		&s.NumeroQuestoes.Tec, &s.TempoRestante, &s.IdUsuario); err != nil {
		return errors.New("Não foi possível obter o simulado.")
	}

	if err := s.getCorrecao(db); err != nil {
		return err
	}

	return nil
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
		return errors.New("Não foi possível finalizar o simulado.")
	}

	return nil

}

func (s *Simulado) Delete(db *sql.DB) error {

	if err := s.getEstado(db); err != nil {
		return err
	}

	if _, err := db.Exec("DELETE FROM simulado WHERE id = $1 AND id_usuario = $2", s.ID, s.IdUsuario); err != nil {
		return errors.New("Não foi possível apagar o simulado.")
	}

	return nil
}

func (s *Simulado) getEstado(db *sql.DB) error {

	s.NumeroQuestoes = &NumeroQuestoes{}
	var user string

	if err := db.QueryRow("SELECT id_usuario, estado, tempo_restante, tempo_limite FROM simulado WHERE id = $1", s.ID).
		Scan(&user, &s.Estado, &s.TempoRestante, &s.TempoLimite); err != nil {
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

	s.Respostas = map[int]int{}
	rows, err := db.Query("SELECT id_questao, resposta FROM questoes_simulado WHERE id_simulado = $1", s.ID)
	if err != nil {
		return errors.New("Não foi possível obter as respostas.")
	}

	for rows.Next() {
		var id int
		var resp int
		err := rows.Scan(&id, &resp)
		s.Respostas[id] = resp

		if err.Error() == `sql: Scan error on column index 1, name "resposta": converting NULL to int is unsupported` {
			s.Respostas[id] = -1
		}

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

func (s *Simulado) getCorrecao(db *sql.DB) error {

	s.Correcao = &Correcao{}

	if err := db.QueryRow("SELECT * FROM correcao WHERE id_simulado = $1", s.ID).
		Scan(&s.Correcao.Brancos.Tot, &s.Correcao.Brancos.Mat, &s.Correcao.Brancos.Fun, &s.Correcao.Brancos.Tec,
			&s.Correcao.Acertos.Tot, &s.Correcao.Acertos.Mat, &s.Correcao.Acertos.Fun, &s.Correcao.Acertos.Tec,
			&s.Correcao.Erros.Tot, &s.Correcao.Erros.Mat, &s.Correcao.Erros.Fun, &s.Correcao.Erros.Tec,
			&s.Correcao.DataFinalizacao); err != nil {
		return errors.New("Não foi possível obter a correção")
	}

	s.Correcao.TempoRealizacao = s.TempoLimite - s.TempoRestante

	return nil

}
