package models

type Comentario struct {
	ID             int    `json:"id"`
	DataPublicacao string `json:"data_publicacao"`
	Texto          string `json:"texto"`
	IdUsuario      string `json:"id_usuario"`
	Sinalizado     int    `json:"sinalizado"`
}
