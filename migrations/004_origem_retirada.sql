-- 004_origem_retirada.sql
-- Data: 2026-09-14
-- Marca de onde o produto e retirado: do armario ou do almoxarifado.
--
-- Isto divide o fluxo do usuario em dois:
--   - armario: retirada direta, com abertura de porta e leitura RFID
--   - almoxarifado: solicitacao, que vira uma lista de separacao para alguem
--     do deposito atender
--
-- Guardado como texto curto em vez de booleano porque podem surgir outras
-- origens (fornecedor externo, outro setor), e assim nao precisa migrar dado.
--
-- Padrao: armario, que e como o sistema funciona hoje.

alter table epiproduto
    add column if not exists origem varchar(16) not null default 'armario';

comment on column epiproduto.origem is 'De onde o produto e retirado: armario ou almoxarifado';
