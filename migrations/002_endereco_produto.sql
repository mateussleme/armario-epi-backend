-- 002_endereco_produto.sql
-- Data: 2026-09-08
-- Adiciona o endereco fisico do produto dentro do armario.
--
-- Formato: letra da linha + numero da coluna, ex: A03 (linha A, coluna 03).
-- Fica gravado ja no formato final, que e como aparece nas telas. No formulario
-- os dois pedacos sao preenchidos separados, para evitar erro de digitacao.
--
-- Aceita nulo porque nem todo produto esta necessariamente dentro do armario
-- (pode estar so no almoxarifado, por exemplo).

alter table epiproduto
    add column if not exists endereco varchar(8);

comment on column epiproduto.endereco is 'Endereco no armario: linha (letra) + coluna (numero), ex: A03';
