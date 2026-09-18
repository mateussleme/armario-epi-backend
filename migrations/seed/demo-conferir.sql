-- demo-conferir.sql
-- Mostra a situacao de cada EPI do usuario, do jeito que o backend calcula.
-- Serve para provar no banco o que a tela esta dizendo.
--
-- Troque 'teste' pelo usuario que estiver demonstrando.

select
    ep.distrib                                             as grupo,
    p.nome                                                 as produto,
    coalesce(l.tipo, 'armario')                            as retira_em,
    case when ep.obrigatorio then 'sim' else 'nao' end     as obrigatorio,
    coalesce(ep.dias_validade::text, 'sem prazo')          as prazo,
    coalesce(to_char(max(r.data), 'DD/MM/YYYY'), 'nunca')  as ultima_retirada,
    case
        when not ep.obrigatorio
            then 'opcional'
        when max(r.data) is null
            then 'COBRAR (nunca retirou)'
        when ep.dias_validade is null
            then 'ok (sem prazo)'
        when extract(day from now() - max(r.data))::int >= ep.dias_validade
            then 'COBRAR (vencido ha '
                 || (extract(day from now() - max(r.data))::int - ep.dias_validade)
                 || ' dias)'
        else 'ok (faltam '
             || (ep.dias_validade - extract(day from now() - max(r.data))::int)
             || ' dias)'
    end                                                    as situacao
from pessoa_epidistrib pe
inner join epidistrib_produtos ep on ep.distrib = pe.distrib
inner join epiproduto p         on p.id = ep.produto
left  join local l              on l.id = p.local
left  join retirada r           on r.pessoa = pe.pessoa and r.produto = ep.produto
where pe.pessoa = 'teste'
group by ep.distrib, p.nome, l.tipo, ep.obrigatorio, ep.dias_validade
order by ep.obrigatorio desc, p.nome;
