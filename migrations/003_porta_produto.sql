-- 003_porta_produto.sql
-- Data: 2026-09-09
-- Adiciona a porta do armario onde o produto fica.
--
-- Sao quatro portas fixas (1 a 4). Guardado como numero: se um dia as portas
-- virarem cadastro proprio (por causa do multiempresa, onde cada armario pode
-- ter uma configuracao), este campo vira a chave estrangeira.
--
-- Aceita nulo porque nem todo produto esta dentro do armario.

alter table epiproduto
    add column if not exists porta smallint;

comment on column epiproduto.porta is 'Porta do armario onde o produto fica (1 a 4)';
