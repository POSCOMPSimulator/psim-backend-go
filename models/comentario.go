package models

import (
	"database/sql"
	"errors"
)

type Comentario struct {
	ID             int    `json:"id"`
	AutorID        string `json:"autor_id"`
	AutorNome      string `json:"autor"`
	QuestaoID      int    `json:"questao_id,omitempty"`
	Texto          string `json:"texto"`
	DataPublicacao string `json:"data_publicacao"`
	Sinalizado     int    `json:"numero_sinalizacoes"`
}

type BatchComentarios struct {
	QuestaoID   int          `json:"-"`
	Comentarios []Comentario `json:"comentarios"`
}

func (bc *BatchComentarios) GetComentariosSinalizados(db *sql.DB) error {

	bc.Comentarios = []Comentario{}
	rows, err := db.Query("SELECT * FROM comentario WHERE sinalizado > 0 ORDER BY sinalizado")
	if err != nil {
		return errors.New("Não foi possível selecionar os comentários.")
	}

	for rows.Next() {
		var comment Comentario
		rows.Scan(&comment.ID, &comment.DataPublicacao, &comment.Texto,
			&comment.AutorID, &bc.QuestaoID, &comment.Sinalizado)
		bc.Comentarios = append(bc.Comentarios, comment)
	}

	return nil

}

func (bc *BatchComentarios) GetComentariosQuestao(db *sql.DB) error {

	query := `
	SELECT id, data_publicacao, texto, id_usuario, nome, id_questao, sinalizado
	FROM comentario
	LEFT JOIN usuario ON usuario.email = comentario.id_usuario
	WHERE id_questao = $1
	ORDER BY data_publicacao DESC
	`

	bc.Comentarios = []Comentario{}
	rows, err := db.Query(query, bc.QuestaoID)
	if err != nil {
		return errors.New("Não foi possível selecionar os comentários.")
	}

	for rows.Next() {
		var comment Comentario
		rows.Scan(&comment.ID, &comment.DataPublicacao, &comment.Texto,
			&comment.AutorID, &comment.AutorNome, &bc.QuestaoID, &comment.Sinalizado)
		bc.Comentarios = append(bc.Comentarios, comment)
	}

	return nil

}

func (c *Comentario) Get(db *sql.DB) error {

	if err := db.QueryRow("SELECT id_usuario FROM comentario WHERE id = $1", c.ID).Scan(&c.AutorID); err != nil {
		return err
	}

	return nil
}

func (c *Comentario) Post(db *sql.DB) error {

	query := "INSERT INTO comentario(data_publicacao, texto, id_usuario, id_questao) VALUES($1, $2, $3, $4)"

	if _, err := db.Exec(query, c.DataPublicacao, c.Texto, c.AutorID, c.QuestaoID); err != nil {
		return errors.New("Não foi possível postar o comentário.")
	}

	return nil
}

func (c *Comentario) Report(db *sql.DB) error {

	if _, err := db.Exec("UPDATE comentario SET sinalizado = sinalizado + 1 WHERE id = $1", c.ID); err != nil {
		return errors.New("Não foi possível reportar o comentário.")
	}

	return nil

}

func (c *Comentario) Clean(db *sql.DB) error {

	if _, err := db.Exec("UPDATE comentario SET sinalizado = 0 WHERE id = $1", c.ID); err != nil {
		return errors.New("Não foi possível limpar o comentário.")
	}

	return nil

}

func (c *Comentario) Delete(db *sql.DB) error {

	if _, err := db.Exec("DELETE FROM comentario WHERE id = $1", c.ID); err != nil {
		return errors.New("Não foi possível deletar o comentário.")
	}

	return nil

}
