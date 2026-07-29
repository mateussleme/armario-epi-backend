package database

import (
	"context"

	"github.com/TopSisErp/epi-backend/internal/utils"
)

func UserData(ctx context.Context, id string) (*User, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, `
		select
			nome,
			imagem,
			case when admin then '1' else '0' end
		from pessoa
		where
			pessoa.id = $1
	`, id)

	name := ""
	image := []byte{}
	admin := "0"
	err = row.Scan(&name, &image, &admin)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:       id,
		Name:     name,
		Admin:    admin == "1",
		ImageUri: utils.BytesToUrl(image),
	}, nil
}

func AllUsers(ctx context.Context) ([]User, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `
		select
			id,
			nome
		from pessoa
		order by id
	`)
	if err != nil {
		return nil, err
	}

	list := []User{}
	defer rows.Close()
	for rows.Next() {
		id := ""
		name := ""

		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		list = append(list, User{Id: id, Name: name})
	}

	return list, nil
}

func UpdateUser(ctx context.Context, id string, name string, image []byte, admin bool) error {
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
			INSERT INTO pessoa (id, nome, imagem, admin)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE 
			SET nome = case when excluded.nome = '' then pessoa.nome else excluded.nome end, 
				imagem = case when excluded.imagem is null then pessoa.imagem else excluded.imagem end,
				admin = excluded.admin;
		`, id, name, image, admin)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteUser(ctx context.Context, id string) error {
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
			DELETE FROM pessoa_epidistrib
			WHERE pessoa_epidistrib.pessoa = $1
		`, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
			DELETE FROM pessoa
			WHERE pessoa.id = $1
		`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func UserProducts(ctx context.Context, id string) ([]string, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `
		select
			epidistrib_produtos.produto
		from pessoa_epidistrib
		inner join epidistrib_produtos on
			epidistrib_produtos.distrib = pessoa_epidistrib.distrib
		where
			pessoa_epidistrib.pessoa = $1
		order by epidistrib_produtos.produto
	`, id)
	if err != nil {
		return nil, err
	}

	products := []string{}
	defer rows.Close()
	for rows.Next() {
		id := ""

		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		products = append(products, id)
	}

	return products, nil
}

func UserAddGroup(ctx context.Context, id string, group string) error {
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
			INSERT INTO pessoa_epidistrib (pessoa, distrib)
			VALUES ($1, $2)
			ON CONFLICT (pessoa, distrib) DO NOTHING
		`, id, group)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func UserRemoveGroup(ctx context.Context, id string, group string) error {
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
			delete from pessoa_epidistrib
			where
				pessoa_epidistrib.pessoa = $1
				and pessoa_epidistrib.distrib = $2
		`, id, group)
	if err != nil {
		return err
	}

	return tx.Commit()
}
