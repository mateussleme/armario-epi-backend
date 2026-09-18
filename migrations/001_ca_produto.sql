-- 001_ca_produto.sql
-- Data: 2026-09-07
-- Adiciona o Certificado de Aprovacao (CA) ao cadastro de produtos.
--
-- O CA e o numero do certificado emitido pelo Ministerio do Trabalho para o EPI,
-- e tem data de vencimento. Sao dados obrigatorios para a ficha de EPI.
--
-- Ambos os campos aceitam nulo porque nem todo item do armario e necessariamente
-- um EPI certificado (o armario pode controlar qualquer coisa, nao so EPI).

alter table epiproduto
    add column if not exists ca varchar(32),
    add column if not exists ca_vencimento date;

comment on column epiproduto.ca is 'Numero do Certificado de Aprovacao (CA) do EPI';
comment on column epiproduto.ca_vencimento is 'Data de vencimento do CA';
