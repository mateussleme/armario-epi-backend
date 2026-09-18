-- 006_validade_uso.sql
-- Data: 2026-09-17
-- Obrigatoriedade e prazo de troca do EPI por grupo (GES).
--
-- Regra que motivou: no GES de eletrica o operador nao pode trabalhar sem luva,
-- e a luva tem que durar 5 dias. A bota tambem e obrigatoria, mas dura 365 dias.
-- Passado o prazo desde a ultima retirada, o sistema obriga a pedir de novo.
--
-- O prazo fica no vinculo grupo x produto, nao no produto: a mesma luva pode
-- durar 5 dias na eletrica e 15 na logistica, porque o desgaste vem da
-- atividade, nao do item.
--
-- obrigatorio e dias_validade sao campos separados de proposito:
--   obrigatorio sem prazo    -> tem que ter, nao vence (capacete, por exemplo)
--   prazo sem obrigatorio    -> sugere a troca, mas nao trava
--   os dois juntos           -> e o caso da luva: trava quando vence

alter table epidistrib_produtos
    add column if not exists obrigatorio    boolean not null default false,
    add column if not exists dias_validade  integer;

comment on column epidistrib_produtos.obrigatorio is 'Item que o GES exige; entra travado no carrinho quando vencido';
comment on column epidistrib_produtos.dias_validade is 'Dias de uso antes da troca, contados da ultima retirada; nulo = nao vence';

-- Historico de retirada. Ate agora o sistema so sabia a quantidade em estoque
-- (epiproduto.quantidade, atualizada pela leitura das tags), nunca quem levou o
-- que e quando. Sem isso nao da para calcular vencimento, que e contado da
-- ultima retirada de cada pessoa.
--
-- Guarda o historico inteiro, uma linha por retirada, em vez de so a ultima
-- data: serve para auditoria de entrega de EPI, que e exigencia da NR-6.
create table if not exists retirada (
    id        bigserial    primary key,
    pessoa    varchar(64)  not null references pessoa(id),
    produto   varchar(64)  not null references epiproduto(id),
    data      timestamptz  not null default now(),
    origem    varchar(16)  not null default 'armario'
);

comment on table  retirada is 'Historico de retirada de EPI por pessoa';
comment on column retirada.origem is 'armario (retirada direta) ou almoxarifado (solicitacao)';

-- A consulta quente e "ultima retirada desta pessoa deste produto".
create index if not exists retirada_pessoa_produto_idx
    on retirada (pessoa, produto, data desc);
