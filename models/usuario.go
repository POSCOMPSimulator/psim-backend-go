package models

import (
	"database/sql"
	"errors"
	"os"
)

type Usuario struct {
	Email        string           `json:"email,omitempty"`
	NivelAcesso  int16            `json:"nivel_acesso"`
	Nome         string           `json:"nome"`
	FotoPerfil   string           `json:"foto_perfil"`
	Estatisticas EstaticasUsuario `json:"estatisticas,omitempty"`
	Completo     bool             `json:"-"`
}

type EstaticasUsuario struct {
	NumSimuladoFinalizado     int                       `json:"num_simulados_finalizados"`
	NumComentariosPublicados  int                       `json:"num_comentarios_publicados"`
	PorcentagemQuestoesFeitas PorcentagemQuestoesFeitas `json:"porcentagem_questoes_feitas"`
}

type PorcentagemQuestoesFeitas struct {
	Geral float32 `json:"geral"`
	Mat   float32 `json:"mat"`
	Fun   float32 `json:"fun"`
	Tec   float32 `json:"tec"`
}

func (u *Usuario) GetDummy() {
	u.Email = os.Getenv("DUMMY_TOKEN")
	u.FotoPerfil = "dummy.png"
	u.Nome = "Dummy User"
	u.NivelAcesso = 32767
}

func (u *Usuario) Create(db *sql.DB) error {
	if _, err := db.Exec("INSERT INTO usuario(email, nome, foto_perfil, nivel_acesso) VALUES($1, $2, $3, $4)", u.Email, u.Nome, u.FotoPerfil, u.NivelAcesso); err != nil {
		return errors.New("Usuário não pode ser criado.")
	}
	return nil
}

func (u *Usuario) Get(db *sql.DB) error {
	if err := db.QueryRow("SELECT email, nome, foto_perfil, nivel_acesso FROM usuario WHERE email=$1", u.Email).Scan(&u.Email, &u.Nome, &u.FotoPerfil, &u.NivelAcesso); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Usuário não encontrado.")
		}
		return err
	}

	if u.Completo {
		_ = db.QueryRow("SELECT count(id_usuario) FROM simulado WHERE id_usuario=$1 AND estado=2", u.Email).Scan(&u.Estatisticas.NumSimuladoFinalizado)
		_ = db.QueryRow("SELECT count(id_usuario) FROM comentario WHERE id_usuario=$1", u.Email).Scan(&u.Estatisticas.NumComentariosPublicados)
		err := u.getQuestoesRealizadas(db)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Usuario) Promote(db *sql.DB) error {
	if err := db.QueryRow("UPDATE usuario SET nivel_acesso = $2 WHERE email = $1 RETURNING email", u.Email, u.NivelAcesso).Scan(&u.Email); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Usuário não encontrado.")
		}
		return err
	}
	return nil
}

func (u *Usuario) Delete(db *sql.DB) error {
	if _, err := db.Exec("DELETE FROM usuario WHERE email=$1", u.Email); err != nil {
		return err
	}
	return nil
}

func (u *Usuario) getQuestoesRealizadas(db *sql.DB) error {
	const queryQuestoesFeitas = `
	SELECT count(distinct(id_questao)) 
	FROM (
		SELECT id_questao, resposta, gabarito, area 
		FROM (
			SELECT id_questao, resposta
			FROM questoes_simulado
			WHERE id_usuario = $1
		) AS q_usuario
		LEFT JOIN questao ON q_usuario.id_questao = questao.id
		WHERE resposta = gabarito
	) AS resp
	GROUP BY area
	`

	rows, err := db.Query(queryQuestoesFeitas, u.Email)
	if err != nil {
		return err
	}

	defer rows.Close()

	var (
		count int
		area  string
		total int
	)

	for rows.Next() {
		rows.Scan(&count, &area)
		switch area {
		case "Matemática":
			u.Estatisticas.PorcentagemQuestoesFeitas.Mat = float32(count)
		case "Fundamentos da Computação":
			u.Estatisticas.PorcentagemQuestoesFeitas.Fun = float32(count)
		case "Tecnologia da Computação":
			u.Estatisticas.PorcentagemQuestoesFeitas.Tec = float32(count)
		}
	}

	u.Estatisticas.PorcentagemQuestoesFeitas.Geral = u.Estatisticas.PorcentagemQuestoesFeitas.Mat +
		u.Estatisticas.PorcentagemQuestoesFeitas.Fun +
		u.Estatisticas.PorcentagemQuestoesFeitas.Tec

	rows, err = db.Query("SELECT count(id) FROM questao GROUP BY area")
	for rows.Next() {
		rows.Scan(&count, &area)
		total += count
		switch area {
		case "Matemática":
			u.Estatisticas.PorcentagemQuestoesFeitas.Mat /= float32(count)
			u.Estatisticas.PorcentagemQuestoesFeitas.Mat *= 100
		case "Fundamentos da Computação":
			u.Estatisticas.PorcentagemQuestoesFeitas.Fun /= float32(count)
			u.Estatisticas.PorcentagemQuestoesFeitas.Fun *= 100
		case "Tecnologia da Computação":
			u.Estatisticas.PorcentagemQuestoesFeitas.Tec /= float32(count)
			u.Estatisticas.PorcentagemQuestoesFeitas.Tec *= 100
		}
	}

	if total > 0 {
		u.Estatisticas.PorcentagemQuestoesFeitas.Geral /= float32(total)
		u.Estatisticas.PorcentagemQuestoesFeitas.Geral *= 100
	}

	return nil
}
