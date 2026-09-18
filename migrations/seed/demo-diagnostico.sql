-- demo-diagnostico.sql
-- Mostra o que existe de verdade no banco, antes de montar a demo.

\echo '=== PRODUTOS CADASTRADOS ==='
select
    p.id,
    p.nome,
    coalesce(l.tipo, '(sem local)') as retira_em,
    p.quantidade
from epiproduto p
left join local l on l.id = p.local
order by p.nome;

\echo ''
\echo '=== VINCULOS ORFAOS (grupo aponta para produto que nao existe) ==='
select distrib as grupo, produto as id_inexistente
from epidistrib_produtos
where produto not in (select id from epiproduto)
order by distrib, produto;

\echo ''
\echo '=== PESSOAS (o admin e quem abre o Gerenciamento) ==='
select
    p.id,
    p.nome,
    case when p.admin then 'ADMIN' else '-' end as perfil,
    string_agg(pe.distrib, ', ' order by pe.distrib) as grupos
from pessoa p
left join pessoa_epidistrib pe on pe.pessoa = p.id
group by p.id, p.nome, p.admin
order by p.admin desc, p.id;

\echo ''
\echo '=== LOCAIS CADASTRADOS ==='
select id, nome, tipo, ativo from local order by tipo, id;
