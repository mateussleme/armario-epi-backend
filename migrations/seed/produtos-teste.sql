-- produtos-teste.sql
-- Popula o cadastro de produtos para testar rolagem de tela e o aviso de CA.
-- ATENCAO: isto e so para o ambiente local de desenvolvimento, nao rodar em producao.
--
-- Os vencimentos foram escolhidos de proposito para cobrir os tres casos:
-- alguns validos, um vencendo em menos de 30 dias e um ja vencido.

insert into epiproduto (id, nome, descricao, quantidade, videourl, ca, ca_vencimento) values
    ('epi4',  'Capacete de segurança',      'Capacete classe B com jugular, para proteção contra impactos na região do crânio.', 18, '', '31469', current_date + interval '18 months'),
    ('epi5',  'Bota de segurança',          'Calçado de segurança com biqueira de composite e solado antiderrapante.',          24, '', '42154', current_date + interval '10 months'),
    ('epi6',  'Protetor facial',            'Protetor facial em policarbonato incolor, para trabalhos com projeção de partículas.', 9, '', '38851', current_date + interval '6 months'),
    ('epi7',  'Máscara PFF2',               'Respirador purificador de ar tipo peça semifacial filtrante para partículas.',      120, '', '38503', current_date + interval '25 days'),
    ('epi8',  'Cinto de segurança',         'Cinturão de segurança tipo paraquedista com talabarte duplo em Y.',                  6, '', '35124', current_date - interval '2 months'),
    ('epi9',  'Avental de raspa',           'Avental de raspa ao couro para proteção do tronco em atividades de solda.',         14, '', '29876', current_date + interval '14 months'),
    ('epi10', 'Manga de raspa',             'Manga de raspa para proteção dos braços contra respingos de solda.',                16, '', '29877', current_date + interval '14 months'),
    ('epi11', 'Óculos de sobreposição',     'Óculos de proteção para uso sobre óculos de grau, com lente incolor.',              30, '', '41260', current_date + interval '20 months'),
    ('epi12', 'Protetor auricular tipo concha', 'Protetor auditivo circum-auricular, atenuação de 19 dB.',                        11, '', '39625', current_date + interval '8 months'),
    ('epi13', 'Luva nitrílica',             'Luva de proteção em nitrilo para manuseio de produtos químicos.',                   45, '', '30112', current_date + interval '3 months'),
    ('epi14', 'Colete refletivo',           'Colete de sinalização com faixas retrorrefletivas.',                                22, '', '33447', current_date + interval '30 months')
on conflict (id) do update
set
    nome = excluded.nome,
    descricao = excluded.descricao,
    quantidade = excluded.quantidade,
    ca = excluded.ca,
    ca_vencimento = excluded.ca_vencimento;
