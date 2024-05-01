-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS questao (
    id serial NOT NULL,
    ano smallint NOT NULL,
    numero smallint NOT NULL,
    area text NOT NULL,
    subarea text,
    alternativa_a text NOT NULL,
    alternativa_b text NOT NULL,
    alternativa_c text NOT NULL,
    alternativa_d text NOT NULL,
    alternativa_e text NOT NULL,
    gabarito int NOT NULL,
    sinalizada boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS simulado (
 id text NOT NULL,
 nome varchar(30),
 estado smallint NOT NULL,
 tempo_limite smallint NOT NULL,
 quant_tot smallint NOT NULL,
 quant_mat smallint NOT NULL,
 quant_fun smallint NOT NULL,
 quant_tec smallint NOT NULL,
 tempo_restante smallint NOT NULL,
 PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS correcao (
    id serial NOT NULL,
    b_total smallint NOT NULL,
    b_mat smallint NOT NULL,
    b_fund smallint NOT NULL,
    b_tec smallint NOT NULL,
    a_total smallint NOT NULL,
    a_mat smallint NOT NULL,
    a_fund smallint NOT NULL,
    a_tec smallint NOT NULL,
    e_total smallint NOT NULL,
    e_mat smallint NOT NULL,
    e_fund smallint NOT NULL,
    e_tec smallint NOT NULL,
    data_finalizacao date NOT NULL,
    id_usuario text NOT NULL,
    id_simulado int NOT NULL,
    PRIMARY KEY (id, id_simulado),
    CONSTRAINT fk_simulado
        FOREIGN KEY (id_simulado)
        REFERENCES simulado(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS questoes_simulado (
    id_simulado int NOT NULL,
    id_usuario text NOT NULL,
    id_questao int NOT NULL,
    resposta int,
    PRIMARY KEY (id_simulado, id_questao),
    CONSTRAINT fk_simulado_q
        FOREIGN KEY (id_simulado)
        REFERENCES simulado(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_questao_sim
        FOREIGN KEY(id_questao)
            REFERENCES questao(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS areas_simulado (
    id_simulado int NOT NULL,
    id_usuario text NOT NULL,
    area text not null,
    PRIMARY KEY (id_simulado, area),
    CONSTRAINT fk_simulado_areas
        FOREIGN KEY (id_simulado)
        REFERENCES simulado(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS anos_simulado (
    id_simulado int NOT NULL,
    id_usuario text NOT NULL,
    ano smallint not null,
    PRIMARY KEY (id_simulado, ano),
    CONSTRAINT fk_simulado_anos
        FOREIGN KEY (id_simulado)
        REFERENCES simulado(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS enunciado_questao (
    id_questao int NOT NULL,
    ordem smallint NOT NULL,
    texto text NOT NULL,
    PRIMARY KEY (id_questao, ordem, texto),
    CONSTRAINT fk_questao_enun
        FOREIGN KEY(id_questao)
            REFERENCES questao(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS imagem_questao (
    id_questao int NOT NULL,
    tipo char(1) NOT NULL,
    url_img text NOT NULL,
    PRIMARY KEY (id_questao, tipo, url_img),
    CONSTRAINT fk_questao_img
        FOREIGN KEY(id_questao)
            REFERENCES questao(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sinalizacao_questao (
    id_questao int NOT NULL,
    msg_err text NOT NULL,
    PRIMARY KEY (id_questao, msg_err),
    CONSTRAINT fk_questao_sin
        FOREIGN KEY(id_questao)
            REFERENCES questao(id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "questao" CASCADE;
DROP TABLE IF EXISTS "simulado" CASCADE;
DROP TABLE IF EXISTS "correcao" CASCADE;
DROP TABLE IF EXISTS "imagem_questao" CASCADE;
DROP TABLE IF EXISTS "questoes_simulado" CASCADE;
DROP TABLE IF EXISTS "areas_simulado" CASCADE;
DROP TABLE IF EXISTS "anos_simulado" CASCADE;
DROP TABLE IF EXISTS "enunciado_questao" CASCADE;
-- +goose StatementEnd
