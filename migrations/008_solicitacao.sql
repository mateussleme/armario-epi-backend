-- 008_solicitacao.sql
-- Data: 2026-09-17
-- Solicitacao de itens do almoxarifado e a separacao (picking).
--
-- Fluxo: o operador monta o carrinho na tela de Solicitacao, confirma, e o
-- pedido cai na Lista de Separacao do gerenciamento. Alguem do almoxarifado
-- abre o pedido, marca item por item conforme separa, e quando termina o pedido
-- fica aguardando o operador buscar.
--
-- Isto vale so para produto de almoxarifado. O que e do armario sai por retirada
-- direta e nao passa por aqui.

create table if not exists solicitacao (
    id              bigserial     primary key,
    pessoa          varchar(64)   not null,
    data            timestamptz   not null default now(),
    -- aberta -> separando -> aguardando_retirada; cancelada sai do fluxo.
    -- O texto fica aberto de proposito: status novo nao pede migration.
    status          varchar(24)   not null default 'aberta',
    -- Quem separou. Fica nulo ate alguem pegar o pedido.
    separador       varchar(64),
    data_separacao  timestamptz
);

comment on table  solicitacao is 'Pedido de EPI do almoxarifado, feito pelo operador na tela de Solicitacao';
comment on column solicitacao.pessoa is 'Id de quem pediu. Sem FK, como em retirada: o pedido e um fato registrado';
comment on column solicitacao.status is 'aberta, separando, aguardando_retirada ou cancelada';

create table if not exists solicitacao_item (
    solicitacao  bigint       not null references solicitacao(id) on delete cascade,
    produto      varchar(64)  not null,
    quantidade   integer      not null default 1,
    -- Item que entrou travado no carrinho por estar vencido no GES. Guardado
    -- aqui porque a regra pode mudar depois, e o pedido tem que continuar
    -- contando a historia de por que aquele item estava la.
    obrigatorio  boolean      not null default false,
    separado     boolean      not null default false,
    primary key (solicitacao, produto)
);

comment on table  solicitacao_item is 'Itens de um pedido, com a marcacao do picking';
comment on column solicitacao_item.obrigatorio is 'Entrou por vencimento de uso, nao por escolha da pessoa';

-- A lista de separacao abre ordenada por data e filtra por status, que sao as
-- duas unicas consultas dessa tela.
create index if not exists solicitacao_status_data_idx
    on solicitacao (status, data desc);

create index if not exists solicitacao_pessoa_idx
    on solicitacao (pessoa, data desc);
