-- demo-retirou-luva.sql
-- Simula a pessoa retirando a luva no armario agora.
--
-- Serve para gravar o antes e depois sem depender do leitor RFID: roda isto,
-- recarrega a tela de Retirada e o aviso da luva sai, porque o prazo passou a
-- contar de hoje. E exatamente o que o TakeConfirm grava quando o armario esta
-- ligado.

do $$
declare
    -- Troque aqui se for demonstrar com outro usuario.
    v_pessoa varchar(64) := 'teste';
    v_luva   varchar(64);
begin
    select id into v_luva from epiproduto where nome ilike '%luva%' order by id limit 1;

    if v_luva is null then
        raise exception 'Nao achei produto com "luva" no nome.';
    end if;

    insert into retirada (pessoa, produto, data, origem)
        values (v_pessoa, v_luva, now(), 'armario');

    raise notice 'Luva (%) retirada agora por %. Recarregue a tela de Retirada.', v_luva, v_pessoa;
end $$;
