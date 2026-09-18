-- regras-ges.sql
-- Regras de obrigatoriedade e prazo de troca, para testar as telas.
-- So para ambiente local de desenvolvimento. Depende de 006_validade_uso.sql.
--
-- Monta o caso que o Clairton descreveu: o operador nao pode trabalhar sem luva,
-- e a luva dura 5 dias. A bota tambem e obrigatoria, mas dura 365 dias.
--
-- A luva fica no armario e a bota no almoxarifado, entao cada uma exercita um
-- fluxo diferente:
--   luva -> aviso na tela de Retirada, antes de abrir a porta
--   bota -> entra travada no carrinho da Solicitacao
--
-- Acha os produtos pelo NOME, nao por id fixo. A versao anterior chutava
-- 'epi1', 'epi5' e afins; como epidistrib_produtos nao valida chave estrangeira,
-- isso criava vinculo apontando para produto inexistente e a tela mostrava
-- "Produto excluido". Aqui so entra o que existe mesmo no cadastro.

do $$
declare
    -- Troque aqui se for demonstrar com outro usuario.
    v_pessoa    varchar(64) := 'teste';
    v_luva      varchar(64);
    v_bota      varchar(64);
    v_armario   varchar(64);
    v_almox     varchar(64);
    v_grupo     varchar(64);
    v_grupos    int := 0;
begin
    -- Limpa o lixo deixado por vinculo de produto que nao existe. Sem isto a
    -- tela fica com card cinza de "Produto excluido".
    delete from epidistrib_produtos
    where produto not in (select id from epiproduto);

    if not exists (select 1 from pessoa where id = v_pessoa) then
        raise exception 'A pessoa % nao existe no cadastro. Cadastre ela ou troque o v_pessoa no topo do bloco.', v_pessoa;
    end if;

    select id into v_luva from epiproduto where nome ilike '%luva%'  order by id limit 1;
    select id into v_bota from epiproduto where nome ilike '%bota%'  order by id limit 1;

    if v_luva is null or v_bota is null then
        raise exception 'Nao achei produto com "luva" ou "bota" no nome. Cadastre os dois antes de rodar a demo.';
    end if;

    select id into v_armario from local where tipo = 'armario'      and ativo order by id limit 1;
    select id into v_almox   from local where tipo = 'almoxarifado' and ativo order by id limit 1;

    -- O cadastro pode nao ter nenhum armario (foi o caso aqui: so existiam dois
    -- almoxarifados). Sem um local do tipo armario nao da para demonstrar a
    -- retirada direta, entao cria. Da no mesmo que cadastrar pela tela de Locais.
    if v_armario is null then
        insert into local (id, nome, tipo, ativo)
            values ('armario01', 'Armário', 'armario', true)
            on conflict (id) do update set ativo = true, tipo = 'armario';
        v_armario := 'armario01';
        raise notice 'Nao havia local do tipo armario. Criei o armario01.';
    end if;

    if v_almox is null then
        insert into local (id, nome, tipo, ativo)
            values ('almoxarifado', 'Almoxarifado', 'almoxarifado', true)
            on conflict (id) do update set ativo = true, tipo = 'almoxarifado';
        v_almox := 'almoxarifado';
        raise notice 'Nao havia local do tipo almoxarifado ativo. Criei o almoxarifado.';
    end if;

    -- Onde cada produto e retirado. Lembrando que almoxarifado aqui quer dizer
    -- "a pessoa pede e alguem separa", nao "o produto so existe la".
    update epiproduto set local = v_armario, origem = 'armario'      where id = v_luva;
    update epiproduto set local = v_almox,   origem = 'almoxarifado' where id = v_bota;

    -- Aplica a regra em todo grupo da pessoa, seja qual for o nome dele.
    for v_grupo in select distinct distrib from pessoa_epidistrib where pessoa = v_pessoa loop
        v_grupos := v_grupos + 1;

        insert into epidistrib_produtos (distrib, produto)
            values (v_grupo, v_luva), (v_grupo, v_bota)
            on conflict (distrib, produto) do nothing;

        update epidistrib_produtos set obrigatorio = true, dias_validade = 5
            where distrib = v_grupo and produto = v_luva;

        update epidistrib_produtos set obrigatorio = true, dias_validade = 365
            where distrib = v_grupo and produto = v_bota;
    end loop;

    if v_grupos = 0 then
        raise exception 'A pessoa % nao esta em nenhum grupo. Vincule ela a um GES na tela de usuarios.', v_pessoa;
    end if;

    -- Historico, para o teste ficar visivel sem esperar 5 dias:
    --   luva retirada ha 9 dias    -> vencida, aparece no aviso do armario
    --   bota retirada ha 400 dias  -> vencida, entra travada no carrinho
    delete from retirada where pessoa = v_pessoa and produto in (v_luva, v_bota);

    insert into retirada (pessoa, produto, data, origem) values
        (v_pessoa, v_luva, now() - interval '9 days',   'armario'),
        (v_pessoa, v_bota, now() - interval '400 days', 'almoxarifado');

    raise notice 'Pronto. luva = % (armario), bota = % (almoxarifado), % grupo(s) de %',
        v_luva, v_bota, v_grupos, v_pessoa;
end $$;
