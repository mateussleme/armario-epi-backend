# Migrations

Os comandos de banco ficam aqui como parte dos fontes, em ordem numerica. Rodar
na sequencia, uma vez cada, num banco que ja tenha o dump inicial restaurado.

```
psql -U postgres -d epi -f migrations/001_ca_produto.sql
psql -U postgres -d epi -f migrations/002_endereco_produto.sql
...
```

| Arquivo | O que faz |
| --- | --- |
| 001_ca_produto.sql | CA e vencimento do CA no produto |
| 002_endereco_produto.sql | Endereco do produto no armario (A01, A02...) |
| 003_porta_produto.sql | Porta do armario onde o produto fica |
| 004_origem_retirada.sql | Origem do produto: armario ou almoxarifado |
| 005_locais.sql | Cadastro de locais; o produto passa a apontar para um local |
| 006_validade_uso.sql | Obrigatoriedade e prazo de troca por GES, e o historico de retirada |
| 007_retirada_historico.sql | Solta as FKs do historico: o registro sobrevive a exclusao do cadastro |
| 008_solicitacao.sql | Solicitacao de itens do almoxarifado e a separacao (picking) |
| 009_quantidade_separada.sql | Quanto foi separado de fato, quando o local nao tem tudo |

Todas sao idempotentes (`if not exists`), entao rodar duas vezes nao quebra.

## seed/

Dados de teste, **so para ambiente local**. Nao rodar em producao.

| Arquivo | O que faz |
| --- | --- |
| produtos-teste.sql | Cadastra EPIs de exemplo |
| vincula-grupos.sql | Liga usuarios e produtos aos GES |
| corrige-acentos.sql | Conserta acentuacao de carga antiga |
| regras-ges.sql | Monta o cenario de obrigatoriedade e prazo (luva 5 dias, bota 365) |
| demo-diagnostico.sql | Mostra produtos, locais, pessoas e vinculos orfaos |
| demo-conferir.sql | Situacao de cada EPI do usuario, do jeito que o backend calcula |
| demo-retirou-luva.sql | Simula uma retirada agora, para testar sem o leitor RFID |
