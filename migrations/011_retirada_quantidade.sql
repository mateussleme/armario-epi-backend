-- 011_retirada_quantidade.sql
-- Data: 2026-09-19
-- Quantidade em cada retirada.
--
-- A ficha de EPI no cadastro do usuario mostra data, produto e quantidade de
-- cada entrega. Ate aqui a tabela guardava uma linha por item, sem quantidade:
-- quem recebeu tres pares de luva aparecia como se tivesse recebido um. Para a
-- ficha isso e erro de registro, nao so de exibicao.
--
-- As linhas antigas ficam com 1, que e o que o armario sempre entregou (uma
-- unidade por retirada).

alter table retirada
    add column if not exists quantidade integer not null default 1;

comment on column retirada.quantidade is 'Quantas unidades a pessoa levou nesta retirada';
