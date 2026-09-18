-- 010_prazo_separacao.sql
-- Data: 2026-09-18
-- Prazo para separar a requisicao.
--
-- Copiado do almoxarifado que ja roda na fabrica: la a TV de requisicoes abertas
-- mostra uma coluna "Separar Ate" e pinta a linha conforme o tempo restante. E o
-- que faz aquela tela funcionar pendurada na parede: de longe se ve laranja e se
-- sabe que tem coisa atrasada, sem ler linha nenhuma.
--
-- O prazo e gravado na requisicao, e nao calculado na hora de exibir, porque
-- mudar o parametro depois nao pode reescrever a historia: requisicao antiga
-- continua sendo cobrada pelo prazo que valia quando ela nasceu.

alter table solicitacao
    add column if not exists separar_ate timestamptz;

comment on column solicitacao.separar_ate is 'Hora limite para separar; gravada na criacao, com o prazo vigente';

-- As que ja existem ganham o prazo contado da propria data, para nao aparecerem
-- sem hora limite na tela.
update solicitacao
set separar_ate = data + interval '15 minutes'
where separar_ate is null;
