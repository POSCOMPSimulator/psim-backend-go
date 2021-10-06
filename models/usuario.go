package models

type Usuario struct {
	GoogleID    string `json:"google_id"`
	NivelAcesso int16  `json:"nivel_acesso"`
}
