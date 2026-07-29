package database

import (
	"context"

	"github.com/TopSisErp/epi-backend/internal/utils"
)

func ItemData(ctx context.Context, id string) (*Item, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, `
		select
			nome,
			descricao,
			imagem,
			videourl
		from epiproduto
		where
			epiproduto.id = $1
	`, id)

	name := ""
	description := ""
	image := []byte{}
	videoUri := ""

	err = row.Scan(&name, &description, &image, &videoUri)
	if err != nil {
		return nil, err
	}

	return &Item{
		Name:        name,
		Description: description,
		ImageUri:    utils.BytesToUrl(image),
		VideoUri:    videoUri,
	}, nil
}

func AllItems(ctx context.Context) ([]Item, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `
		select
			id,
			nome,
			descricao
		from epiproduto
		order by id
	`)
	if err != nil {
		return nil, err
	}

	list := []Item{}
	defer rows.Close()
	for rows.Next() {
		id := ""
		name := ""
		description := ""

		err := rows.Scan(&id, &name, &description)
		if err != nil {
			return nil, err
		}

		list = append(list, Item{Id: id, Name: name, Description: description})
	}

	return list, nil
}

func UpdateItem(ctx context.Context, id string, name string, description string, image []byte, videoUri string) error {
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
			INSERT INTO epiproduto (id, nome, descricao, imagem, videourl, quantidade)
			VALUES ($1, $2, $3, $4, $5, 0)
			ON CONFLICT (id) DO UPDATE 
			SET nome = case when excluded.nome = '' then epiproduto.nome else excluded.nome end,
				descricao = case when excluded.descricao = '' then epiproduto.descricao else excluded.descricao end,
				imagem = case when excluded.imagem is null then epiproduto.imagem else excluded.imagem end,
				videourl = case when excluded.videourl = '' then epiproduto.videourl else excluded.videourl end;
		`, id, name, description, image, videoUri)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteItem(ctx context.Context, id string) error {
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
			WHERE epidistrib_produtos.produto = $1
		`, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
			DELETE FROM epiproduto
			WHERE epiproduto.id = $1
		`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
