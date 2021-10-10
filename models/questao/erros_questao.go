package questao

import (
	"database/sql"
	"errors"
	"strconv"
)

type ErrosQuestao struct {
	ID    int      `json:"-"`
	Erros []string `json:"erros"`
}

type MensagemErro struct {
	ID   int    `json:"id_questao"`
	Erro string `json:"mensagem_erro"`
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

	queryString := "DELETE FROM sinalizacao_questao WHERE id_questao = $1 AND msg_err IN ("
	msgs := []interface{}{eq.ID}

	seps := map[bool]string{true: ",$", false: "$"}

	if len(eq.Erros) == 0 {
		return nil
	}

	for i, v := range eq.Erros {
		queryString += seps[i > 0] + strconv.Itoa(i+2)
		msgs = append(msgs, v)
	}
	queryString += ")"

	if _, err := db.Exec(queryString, msgs...); err != nil {
		return errors.New("Não foi possível resolver os erros.")
	}

	othereq := ErrosQuestao{ID: eq.ID, Erros: []string{}}
	othereq.Get(db)

	if len(othereq.Erros) == 0 {
		db.Exec("UPDATE questao SET sinalizada = false WHERE id = $1", eq.ID)
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
