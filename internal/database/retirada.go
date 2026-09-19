package database

import (
	"context"
	"database/sql"
	"time"
)

// Marca que a pessoa levou o produto agora. E daqui que sai a contagem do prazo
// de troca: o vencimento e sempre contado da ultima retirada de cada pessoa,
// nao de uma data fixa do produto.
func RegisterRetirada(ctx context.Context, pessoa string, produto string, origem string) error {
	conn, err := Connection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if origem != "almoxarifado" {
		origem = "armario"
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
			INSERT INTO retirada (pessoa, produto, origem)
			VALUES ($1, $2, $3)
		`, pessoa, produto, origem)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Os itens que a pessoa e obrigada a pedir agora.
//
// Entra na lista o produto que e obrigatorio em algum GES dela e que:
//   - ela nunca retirou, ou
//   - passou do prazo desde a ultima retirada
//
// Quando o mesmo produto aparece em mais de um GES da pessoa, vale a regra mais
// apertada: basta um GES marcar como obrigatorio, e o prazo e o menor entre
// eles. O min do SQL ignora nulo, entao um GES sem prazo nao afrouxa o outro.
//
// Produto que ja esta num pedido em aberto sai da lista: a pessoa fez a parte
// dela e esta esperando o almoxarifado separar, entao cobrar de novo so faria
// ela abrir pedido em cima de pedido. Volta a ser cobrado se o pedido for
// cancelado.
//
// Enquanto nao existir a entrega, um pedido separado fica em aguardando
// retirada para sempre e o item nao volta a ser cobrado. E o prazo de validade
// do pedido (as 48h que o Clairton mencionou) que vai fechar esse buraco.
func UserRequiredProducts(ctx context.Context, pessoa string) ([]RequiredProduct, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `
		select
			ep.produto,
			min(ep.dias_validade) as dias_validade,
			max(r.data) as ultima_retirada
		from pessoa_epidistrib pe
		inner join epidistrib_produtos ep on
			ep.distrib = pe.distrib
		left join retirada r on
			r.pessoa = pe.pessoa
			and r.produto = ep.produto
		where
			pe.pessoa = $1
			and not exists (
				select 1
				from solicitacao s
				inner join solicitacao_item si on si.solicitacao = s.id
				where
					s.pessoa = pe.pessoa
					and si.produto = ep.produto
					and s.status <> 'cancelada'
			)
		group by ep.produto
		having bool_or(ep.obrigatorio)
		order by ep.produto
	`, pessoa)
	if err != nil {
		return nil, err
	}

	required := []RequiredProduct{}
	defer rows.Close()
	for rows.Next() {
		produto := ""
		dias := sql.NullInt64{}
		ultima := sql.NullTime{}

		err := rows.Scan(&produto, &dias, &ultima)
		if err != nil {
			return nil, err
		}

		item := RequiredProduct{
			Produto:      produto,
			DiasDesdeUso: -1,
		}
		if dias.Valid {
			item.DiasValidade = int(dias.Int64)
		}

		// Nunca retirou: e obrigatorio e a pessoa nao tem. Entra na lista com ou
		// sem prazo definido.
		if !ultima.Valid {
			required = append(required, item)
			continue
		}

		item.UltimaRetirada = ultima.Time.Format(time.RFC3339)
		item.DiasDesdeUso = int(time.Since(ultima.Time).Hours() / 24)

		// Obrigatorio sem prazo nao vence: ja retirou uma vez, esta cumprido.
		if !dias.Valid {
			continue
		}

		if item.DiasDesdeUso >= item.DiasValidade {
			required = append(required, item)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return required, nil
}

// Historico de retirada da pessoa, da mais recente para a mais antiga.
// Usado pela auditoria de entrega de EPI.
func UserRetiradas(ctx context.Context, pessoa string, limit int) ([]Retirada, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if limit <= 0 {
		limit = 100
	}

	rows, err := conn.QueryContext(ctx, `
		select
			r.produto,
			coalesce(p.nome, ''),
			r.quantidade,
			r.data,
			r.origem
		from retirada r
		left join epiproduto p on p.id = r.produto
		where
			r.pessoa = $1
		order by r.data desc
		limit $2
	`, pessoa, limit)
	if err != nil {
		return nil, err
	}

	list := []Retirada{}
	defer rows.Close()
	for rows.Next() {
		item := Retirada{Pessoa: pessoa}
		data := time.Time{}

		err := rows.Scan(&item.Produto, &item.ProdutoNome, &item.Quantidade, &data, &item.Origem)
		if err != nil {
			return nil, err
		}

		item.Data = data.Format(time.RFC3339)
		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
