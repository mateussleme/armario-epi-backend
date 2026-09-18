-- Corrige os acentos dos produtos de teste.
-- O problema: o psql no Windows leu o arquivo UTF-8 usando outra codificacao,
-- entao os acentos viraram sequencias estranhas (nitrÃ¡lica em vez de nitrilica).
--
-- Antes de rodar este arquivo, executar no psql:
--   \encoding UTF8

update epiproduto set nome = 'Capacete de segurança',           descricao = 'Capacete classe B com jugular, para proteção contra impactos na região do crânio.' where id = 'epi4';
update epiproduto set nome = 'Bota de segurança',               descricao = 'Calçado de segurança com biqueira de composite e solado antiderrapante.' where id = 'epi5';
update epiproduto set nome = 'Protetor facial',                 descricao = 'Protetor facial em policarbonato incolor, para trabalhos com projeção de partículas.' where id = 'epi6';
update epiproduto set nome = 'Máscara PFF2',                    descricao = 'Respirador purificador de ar tipo peça semifacial filtrante para partículas.' where id = 'epi7';
update epiproduto set nome = 'Cinto de segurança',              descricao = 'Cinturão de segurança tipo paraquedista com talabarte duplo em Y.' where id = 'epi8';
update epiproduto set nome = 'Avental de raspa',                descricao = 'Avental de raspa ao couro para proteção do tronco em atividades de solda.' where id = 'epi9';
update epiproduto set nome = 'Manga de raspa',                  descricao = 'Manga de raspa para proteção dos braços contra respingos de solda.' where id = 'epi10';
update epiproduto set nome = 'Óculos de sobreposição',          descricao = 'Óculos de proteção para uso sobre óculos de grau, com lente incolor.' where id = 'epi11';
update epiproduto set nome = 'Protetor auricular tipo concha',  descricao = 'Protetor auditivo circum-auricular, atenuação de 19 dB.' where id = 'epi12';
update epiproduto set nome = 'Luva nitrílica',                  descricao = 'Luva de proteção em nitrilo para manuseio de produtos químicos.' where id = 'epi13';
update epiproduto set nome = 'Colete refletivo',                descricao = 'Colete de sinalização com faixas retrorrefletivas.' where id = 'epi14';
