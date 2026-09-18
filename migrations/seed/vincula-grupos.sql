-- vincula-grupos.sql
-- Popula os vinculos para a tela de Retirada ter o que mostrar.
-- So para ambiente local de desenvolvimento.
--
-- A tela de Retirada mostra os EPIs dos grupos a que o usuario pertence
-- (ver GetAvailableItems no front e UserProducts no backend). Sem vinculo,
-- a tela aparece vazia mesmo com produtos cadastrados.

-- Garante que os grupos existem
insert into epidistrib (id, nome) values
    ('soldagem',  'Grupo de Soldagem'),
    ('eletrica',  'Grupo de Elétrica'),
    ('logistica', 'Grupo de Logística')
on conflict (id) do update set nome = excluded.nome;

-- Produtos de cada grupo
insert into epidistrib_produtos (distrib, produto) values
    -- soldagem: protecao para solda
    ('soldagem', 'epi1'),   -- luva
    ('soldagem', 'epi2'),   -- oculos
    ('soldagem', 'epi4'),   -- capacete
    ('soldagem', 'epi6'),   -- protetor facial
    ('soldagem', 'epi9'),   -- avental de raspa
    ('soldagem', 'epi10'),  -- manga de raspa
    ('soldagem', 'epi5'),   -- bota
    -- eletrica
    ('eletrica', 'epi2'),
    ('eletrica', 'epi4'),
    ('eletrica', 'epi5'),
    ('eletrica', 'epi13'),  -- luva nitrilica
    -- logistica
    ('logistica', 'epi3'),  -- protetor auricular
    ('logistica', 'epi5'),
    ('logistica', 'epi14')  -- colete refletivo
on conflict do nothing;

-- Vincula os usuarios de teste aos grupos
insert into pessoa_epidistrib (pessoa, distrib) values
    ('teste',    'soldagem'),
    ('clairton', 'soldagem'),
    ('clairton', 'eletrica')
on conflict do nothing;
