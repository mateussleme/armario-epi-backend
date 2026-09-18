-- 007_retirada_historico.sql
-- Data: 2026-09-17
-- Solta as chaves estrangeiras da tabela retirada.
--
-- A 006 criou retirada apontando para pessoa(id) e epiproduto(id). Parecia certo,
-- mas inverte a prioridade: com a FK, excluir um produto ou uma pessoa que ja
-- tem entrega registrada passa a ser recusado pelo banco, e a tela so diz "nao
-- foi possivel excluir" sem explicar.
--
-- O historico de entrega de EPI e documento de auditoria (NR-6): ele registra um
-- fato que aconteceu, e nao deixa de ter acontecido porque alguem apagou o
-- cadastro depois. Entao o certo e o contrario do que estava: o cadastro pode
-- sair, o registro fica, guardando o id de quem levou e do que foi levado.
--
-- O preco e que o historico pode citar um id que nao existe mais. Isso e aceito
-- de proposito: quem le o historico quer saber o que foi entregue naquele dia,
-- nao o estado atual do cadastro.

alter table retirada drop constraint if exists retirada_pessoa_fkey;
alter table retirada drop constraint if exists retirada_produto_fkey;

comment on column retirada.pessoa is 'Id de quem retirou. Sem FK: o registro sobrevive a exclusao do cadastro';
comment on column retirada.produto is 'Id do que foi retirado. Sem FK: o registro sobrevive a exclusao do cadastro';
