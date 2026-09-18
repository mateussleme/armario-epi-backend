package database

import (
	"context"
	"database/sql"
)

func LocalData(ctx context.Context, id string) (Local, error) {
	var local Local
	var empresa sql.NullString

	err := Db.QueryRowContext(ctx, `
		select
			id,
			nome,
			tipo,
			empresa,
			ativo
		from local
		where
			id = $1
	`, id).Scan(&local.Id, &local.Nome, &local.Tipo, &empresa, &local.Ativo)
	if err != nil {
		return local, err
	}

	if empresa.Valid {
		local.Empresa = empresa.String
	}

	return local, nil
}

// Traz todos, inclusive os desativados: a tela de cadastro precisa mostrar para
// poder reativar. Quem so quer os disponiveis filtra por ativo.
func AllLocais(ctx context.Context) ([]Local, error) {
	rows, err := Db.QueryContext(ctx, `
		select
			id,
			nome,
			tipo,
			empresa,
			ativo
		from local
		order by nome
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	locais := []Local{}
	for rows.Next() {
		var local Local
		var empresa sql.NullString

		if err := rows.Scan(&local.Id, &local.Nome, &local.Tipo, &empresa, &local.Ativo); err != nil {
			return nil, err
		}

		if empresa.Valid {
			local.Empresa = empresa.String
		}

		locais = append(locais, local)
	}

	return locais, rows.Err()
}

func UpdateLocal(ctx context.Context, id string, nome string, tipo string, ativo bool) error {
	_, err := Db.ExecContext(ctx, `
		insert into local (id, nome, tipo, ativo)
		values ($1, $2, $3, $4)
		on conflict (id) do update
		set
			nome = excluded.nome,
			tipo = excluded.tipo,
			ativo = excluded.ativo
	`, id, nome, tipo, ativo)

	return err
}

// Excluir de verdade so funciona se nenhum produto apontar para o local. Quando
// houver produto vinculado, o banco recusa por causa da chave estrangeira, e a
// tela orienta a desativar em vez de excluir.
func DeleteLocal(ctx context.Context, id string) error {
	_, err := Db.ExecContext(ctx, `
		delete from local
		where
			id = $1
	`, id)

	return err
}
