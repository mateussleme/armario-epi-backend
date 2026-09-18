package database

import (
	"context"
	"database/sql"

	"github.com/TopSisErp/epi-backend/internal/utils"
)

func ItemData(ctx context.Context, id string) (Item, error) {
	var item Item
	var image []byte
	var videoUri sql.NullString
	var ca sql.NullString
	var caVencimento sql.NullTime
	var endereco sql.NullString
	var porta sql.NullInt16
	var local sql.NullString

	err := Db.QueryRowContext(ctx, `
		select
			id,
			nome,
			descricao,
			imagem,
			videourl,
			ca,
			ca_vencimento,
			endereco,
			porta,
			local
		from epiproduto
		where
			id = $1
	`, id).Scan(&item.Id, &item.Name, &item.Description, &image, &videoUri, &ca, &caVencimento, &endereco, &porta, &local)
	if err != nil {
		return item, err
	}

	item.ImageUri = utils.BytesToUrl(image)
	if videoUri.Valid {
		item.VideoUri = videoUri.String
	}
	if ca.Valid {
		item.Ca = ca.String
	}
	// A data vai para o front no formato ISO (YYYY-MM-DD), que e o que o input
	// type="date" do HTML espera.
	if caVencimento.Valid {
		item.CaVencimento = caVencimento.Time.Format("2006-01-02")
	}
	if endereco.Valid {
		item.Endereco = endereco.String
	}
	if porta.Valid {
		item.Porta = int(porta.Int16)
	}
	if local.Valid {
		item.Local = local.String
	}

	return item, nil
}

func AllItems(ctx context.Context) ([]Item, error) {
	rows, err := Db.QueryContext(ctx, `
		select
			id,
			nome,
			descricao,
			ca,
			ca_vencimento,
			endereco,
			porta,
			local
		from epiproduto
		order by nome
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item
		var ca sql.NullString
		var caVencimento sql.NullTime
		var endereco sql.NullString
		var porta sql.NullInt16
		var local sql.NullString

		if err := rows.Scan(&item.Id, &item.Name, &item.Description, &ca, &caVencimento, &endereco, &porta, &local); err != nil {
			return nil, err
		}

		if ca.Valid {
			item.Ca = ca.String
		}
		if caVencimento.Valid {
			item.CaVencimento = caVencimento.Time.Format("2006-01-02")
		}
		if endereco.Valid {
			item.Endereco = endereco.String
		}
		if porta.Valid {
			item.Porta = int(porta.Int16)
		}
		if local.Valid {
			item.Local = local.String
		}

		items = append(items, item)
	}

	return items, rows.Err()
}

func UpdateItem(ctx context.Context, id string, name string, description string, image []byte, videoUri string, ca string, caVencimento string, endereco string, porta int, local string) error {
	// Campos vazios viram nulo no banco, para nao gravar valor onde a informacao
	// simplesmente nao foi preenchida.
	var caValue any
	if ca != "" {
		caValue = ca
	}

	var caVencimentoValue any
	if caVencimento != "" {
		caVencimentoValue = caVencimento
	}

	var enderecoValue any
	if endereco != "" {
		enderecoValue = endereco
	}

	var portaValue any
	if porta > 0 {
		portaValue = porta
	}

	// Local vazio fica nulo. A chave estrangeira recusaria uma string vazia, que
	// nao corresponde a nenhum local cadastrado.
	var localValue any
	if local != "" {
		localValue = local
	}

	// A imagem so e sobrescrita quando vem uma nova; senao mantem a que ja existe.
	_, err := Db.ExecContext(ctx, `
		insert into epiproduto (id, nome, descricao, imagem, videourl, quantidade, ca, ca_vencimento, endereco, porta, local)
		values ($1, $2, $3, $4, $5, 0, $6, $7, $8, $9, $10)
		on conflict (id) do update
		set
			nome = excluded.nome,
			descricao = excluded.descricao,
			imagem = case when $4::bytea is null then epiproduto.imagem else excluded.imagem end,
			videourl = excluded.videourl,
			ca = excluded.ca,
			ca_vencimento = excluded.ca_vencimento,
			endereco = excluded.endereco,
			porta = excluded.porta,
			local = excluded.local
	`, id, name, description, image, videoUri, caValue, caVencimentoValue, enderecoValue, portaValue, localValue)

	return err
}

// Exclui o produto e o que aponta para ele.
//
// Antes so apagava de epiproduto, e o vinculo com os grupos ficava para tras
// apontando para o nada. O efeito aparecia na tela do operador: o item virava um
// card cinza de "Produto excluido" e, se estivesse marcado como obrigatorio no
// GES, virava uma cobranca que ninguem conseguia cumprir, porque o produto nao
// existia mais para ser retirado.
//
// O historico de retirada nao e tocado: ele registra entrega que aconteceu, e
// vale como auditoria mesmo depois de o cadastro sair (ver 007).
func DeleteItem(ctx context.Context, id string) error {
	tx, err := Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Vinculo com os grupos: e o que virava cobranca impossivel.
	_, err = tx.ExecContext(ctx, `
		delete from epidistrib_produtos
		where
			epidistrib_produtos.produto = $1
	`, id)
	if err != nil {
		return err
	}

	// Tags RFID do produto. Sem isto a antena continuaria lendo etiqueta de um
	// produto que o cadastro nao conhece, e a leitura cairia em tag desconhecida.
	_, err = tx.ExecContext(ctx, `
		delete from epitag
		where
			epitag.produto = $1
	`, id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		delete from epiproduto
		where
			id = $1
	`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
