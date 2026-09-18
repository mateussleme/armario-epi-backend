-- 009_quantidade_separada.sql
-- Data: 2026-09-17
-- Quanto foi realmente separado de cada item do pedido.
--
-- Ate agora o picking era so um sim ou nao. Mas o almoxarifado nem sempre tem o
-- que foi pedido: se o cara pediu 3 pares e o local tem 2, ele separa 2 e o
-- pedido precisa registrar isso, senao o sistema acha que entregou 3 e o saldo
-- e a ficha de EPI ficam mentindo.
--
-- Nulo enquanto o item nao foi separado. Quando marca, grava o que saiu de fato.

alter table solicitacao_item
    add column if not exists quantidade_separada integer;

comment on column solicitacao_item.quantidade_separada is 'Quanto saiu de fato; pode ser menor que quantidade quando falta saldo';
