package models

type Simulado struct {
	ID            int      `json:"id"`
	Nome          string   `json:"nome"`
	Estado        int      `json:"estado"`
	TempoLimite   int      `json:"tempo_limite"`
	QuantTot      int      `json:"quant_tot"`
	QuantMat      int      `json:"quant_mat"`
	QuantFun      int      `json:"quant_fun"`
	QuantTec      int      `json:"quant_tec"`
	TempoRestante int      `json:"tempo_restante"`
	IdUsuario     string   `json:"id_usuario"`
	Correcao      Correcao `json:"correcao"`
	Anos          []int    `json:"anos"`
	Areas         []string `json:"areas"`
}

type Correcao struct {
}
