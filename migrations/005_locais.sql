-- 005_locais.sql
-- Data: 2026-09-16
-- Cadastro de locais de retirada.
--
-- Antes o local era um select fixo com duas opcoes (armario e almoxarifado).
-- Vira cadastro porque pode existir mais de um armario e mais de um
-- almoxarifado, cada um numa empresa (armario 01 CTBA, almoxarifado SJP...).
--
-- Dois conceitos diferentes convivem aqui:
--   tipo  -> define o fluxo: armario e retirada direta, almoxarifado passa por
--            solicitacao e separacao
--   local -> onde fisicamente, para quando houver mais de um
--
-- O campo empresa ja fica previsto, mesmo sem uso hoje: quando o multiempresa
-- chegar, e so preencher, sem migrar dado.

create table if not exists local (
    id        varchar(64)  primary key,
    nome      varchar(128) not null,
    tipo      varchar(16)  not null,
    empresa   varchar(64),
    ativo     boolean      not null default true
);

comment on table  local is 'Locais de onde os produtos sao retirados';
comment on column local.tipo is 'armario ou almoxarifado; define o fluxo de retirada';
comment on column local.empresa is 'Reservado para o multiempresa; nulo por enquanto';
comment on column local.ativo is 'Local desativado nao aparece para escolha, mas o historico continua valendo';

-- Os dois locais que existem hoje.
insert into local (id, nome, tipo) values
    ('armario01',    'Armário',     'armario'),
    ('almoxarifado', 'Almoxarifado', 'almoxarifado')
on conflict (id) do nothing;

-- O produto passa a apontar para o local em vez de guardar o tipo direto.
alter table epiproduto
    add column if not exists local varchar(64) references local(id);

comment on column epiproduto.local is 'De onde o produto e retirado (referencia local.id)';

-- Migra o que ja estava gravado na coluna origem.
update epiproduto set local = 'armario01'    where local is null and (origem is null or origem = 'armario');
update epiproduto set local = 'almoxarifado' where local is null and origem = 'almoxarifado';

-- A coluna origem fica por enquanto, para nao quebrar nada que ainda leia dela.
-- Quando tudo estiver usando local, da para remover com:
--   alter table epiproduto drop column origem;
