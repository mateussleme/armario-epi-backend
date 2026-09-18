package database

import (
	"context"
	"database/sql"
)

func GroupData(ctx context.Context, id string) (*Group, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, `
		select
			nome
		from epidistrib
		where
			epidistrib.id = $1
	`, id)

	name := ""
	err = row.Scan(&name)
	if err != nil {
		return nil, err
	}

	return &Group{
		Id:   id,
		Name: name,
	}, nil
}

func AllGroups(ctx context.Context) ([]Group, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `
		select
			id,
			nome
		from epidistrib
		order by id
	`)
	if err != nil {
		return nil, err
	}

	list := []Group{}
	defer rows.Close()
	for rows.Next() {
		id := ""
		name := ""

		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		list = append(list, Group{Id: id, Name: name})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func UpdateGroup(ctx context.Context, id string, name string) error {
	conn, err := Connection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
			INSERT INTO epidistrib (id, nome)
			VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE 
			SET nome = case when excluded.nome = '' then epidistrib.nome else excluded.nome end
		`, id, name)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteGroup(ctx context.Context, id string) error {
	conn, err := Connection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
			DELETE FROM epidistrib_produtos
			WHERE epidistrib_produtos.distrib = $1
		`, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
			DELETE FROM pessoa_epidistrib
			WHERE pessoa_epidistrib.distrib = $1
		`, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
			DELETE FROM epidistrib
			WHERE epidistrib.id = $1
		`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GroupAddProduct(ctx context.Context, id string, product string) error {
	conn, err := Connection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
			INSERT INTO epidistrib_produtos (distrib, produto)
			VALUES ($1, $2)
			ON CONFLICT (distrib, produto) DO NOTHING
		`, id, product)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GroupRemoveProduct(ctx context.Context, id string, product string) error {
	conn, err := Connection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
			delete from epidistrib_produtos
			where
				epidistrib_produtos.distrib = $1
				and epidistrib_produtos.produto = $2
		`, id, product)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Os produtos do grupo com a regra de cada um. Antes nao existia rota que
// devolvesse isso, entao a tela de grupo abria sempre vazia e o que ja estava
// gravado parecia ter sumido.
func GroupProducts(ctx context.Context, groupId string) ([]GroupProduct, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `
		select
			produto,
			obrigatorio,
			dias_validade
		from epidistrib_produtos
		where
			epidistrib_produtos.distrib = $1
		order by produto
	`, groupId)
	if err != nil {
		return nil, err
	}

	products := []GroupProduct{}
	defer rows.Close()
	for rows.Next() {
		item := GroupProduct{}
		dias := sql.NullInt64{}

		err := rows.Scan(&item.Produto, &item.Obrigatorio, &dias)
		if err != nil {
			return nil, err
		}

		if dias.Valid {
			value := int(dias.Int64)
			item.DiasValidade = &value
		}

		products = append(products, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// Grava a regra do item dentro do grupo. Mesma rota do adicionar: se o vinculo
// ainda nao existe, cria; se existe, so atualiza a regra. Evita a tela ter que
// saber se esta criando ou editando o vinculo.
//
// dias nulo significa que o item nao vence.
func GroupSetProduct(ctx context.Context, id string, product string, obrigatorio bool, dias *int) error {
	conn, err := Connection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
			INSERT INTO epidistrib_produtos (distrib, produto, obrigatorio, dias_validade)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (distrib, produto) DO UPDATE
			SET obrigatorio = excluded.obrigatorio,
				dias_validade = excluded.dias_validade
		`, id, product, obrigatorio, dias)
	if err != nil {
		return err
	}

	return tx.Commit()
}
