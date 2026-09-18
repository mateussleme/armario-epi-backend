package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// Cancelamento recusado: ou o pedido nao e da pessoa, ou a separacao ja comecou.
var ErrCancelamentoNaoPermitido = errors.New("cancelamento nao permitido")

// Minutos para separar uma requisicao, contados da criacao. O almoxarifado da
// fabrica trabalha com quinze; aqui fica em variavel de ambiente porque o numero
// depende do tamanho do deposito e de quanta gente separa, e trocar isso nao
// pode exigir recompilar.
func prazoSeparacao() time.Duration {
	minutos, err := strconv.Atoi(os.Getenv("PRAZO_SEPARACAO_MINUTOS"))
	if err != nil || minutos <= 0 {
		minutos = 15
	}
	return time.Duration(minutos) * time.Minute
}

// Item como ele chega da tela na hora de criar o pedido.
type NovaSolicitacaoItem struct {
	Produto     string
	Quantidade  int
	Obrigatorio bool
}

// Grava o pedido inteiro de uma vez. Cabecalho e itens na mesma transacao: um
// pedido sem itens nao serve para nada e apareceria na lista de separacao como
// uma linha vazia.
func CreateSolicitacao(ctx context.Context, pessoa string, itens []NovaSolicitacaoItem) (int64, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	id := int64(0)
	// O prazo e resolvido aqui e gravado junto: assim a requisicao carrega o
	// prazo que valia quando ela nasceu, e mudar o parametro depois nao muda o
	// passado.
	err = tx.QueryRowContext(ctx, `
		INSERT INTO solicitacao (pessoa, status, separar_ate)
		VALUES ($1, $2, now() + $3::interval)
		RETURNING id
	`, pessoa, StatusAberta, fmt.Sprintf("%d minutes", int(prazoSeparacao().Minutes()))).Scan(&id)
	if err != nil {
		return 0, err
	}

	for _, item := range itens {
		quantidade := item.Quantidade
		if quantidade <= 0 {
			quantidade = 1
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO solicitacao_item (solicitacao, produto, quantidade, obrigatorio)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (solicitacao, produto) DO UPDATE
			SET quantidade = solicitacao_item.quantidade + excluded.quantidade
		`, id, item.Produto, quantidade, item.Obrigatorio)
		if err != nil {
			return 0, err
		}
	}

	return id, tx.Commit()
}

// A Lista de Separacao. status vazio traz tudo que ainda esta em aberto, que e
// como a tela abre; cancelada nunca entra no padrao, porque sair da lista e o
// proposito dela.
func AllSolicitacoes(ctx context.Context, status string) ([]Solicitacao, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `
		select
			s.id,
			s.pessoa,
			coalesce(p.nome, s.pessoa) as pessoa_nome,
			s.data,
			s.separar_ate,
			s.status,
			coalesce(s.separador, ''),
			s.data_separacao,
			count(i.produto) as total,
			count(i.produto) filter (where i.separado) as separados
		from solicitacao s
		left join pessoa p on p.id = s.pessoa
		left join solicitacao_item i on i.solicitacao = s.id
		where
			-- Sem filtro, a lista de separacao traz so o que ainda da trabalho.
			-- Cancelada e entregue ja sairam do fluxo.
			($1 = '' and s.status <> $2 and s.status <> $3)
			or s.status = $1
		group by s.id, s.pessoa, p.nome, s.data, s.separar_ate, s.status, s.separador, s.data_separacao
		order by s.separar_ate
	`, status, StatusCancelada, StatusEntregue)
	if err != nil {
		return nil, err
	}

	list := []Solicitacao{}
	defer rows.Close()
	for rows.Next() {
		item := Solicitacao{}
		data := time.Time{}
		separarAte := sql.NullTime{}
		dataSeparacao := sql.NullTime{}

		err := rows.Scan(
			&item.Id, &item.Pessoa, &item.PessoaNome, &data, &separarAte, &item.Status,
			&item.Separador, &dataSeparacao, &item.TotalItens, &item.ItensSeparados,
		)
		if err != nil {
			return nil, err
		}

		item.Data = data.Format(time.RFC3339)
		if separarAte.Valid {
			item.SepararAte = separarAte.Time.Format(time.RFC3339)
		}
		if dataSeparacao.Valid {
			item.DataSeparacao = dataSeparacao.Time.Format(time.RFC3339)
		}

		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func SolicitacaoData(ctx context.Context, id int64) (*Solicitacao, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	item := Solicitacao{}
	data := time.Time{}
	separarAte := sql.NullTime{}
	dataSeparacao := sql.NullTime{}

	err = conn.QueryRowContext(ctx, `
		select
			s.id,
			s.pessoa,
			coalesce(p.nome, s.pessoa) as pessoa_nome,
			s.data,
			s.separar_ate,
			s.status,
			coalesce(s.separador, ''),
			s.data_separacao,
			count(i.produto) as total,
			count(i.produto) filter (where i.separado) as separados
		from solicitacao s
		left join pessoa p on p.id = s.pessoa
		left join solicitacao_item i on i.solicitacao = s.id
		where s.id = $1
		group by s.id, s.pessoa, p.nome, s.data, s.separar_ate, s.status, s.separador, s.data_separacao
	`, id).Scan(
		&item.Id, &item.Pessoa, &item.PessoaNome, &data, &separarAte, &item.Status,
		&item.Separador, &dataSeparacao, &item.TotalItens, &item.ItensSeparados,
	)
	if err != nil {
		return nil, err
	}

	item.Data = data.Format(time.RFC3339)
	if separarAte.Valid {
		item.SepararAte = separarAte.Time.Format(time.RFC3339)
	}
	if dataSeparacao.Valid {
		item.DataSeparacao = dataSeparacao.Time.Format(time.RFC3339)
	}

	return &item, nil
}

// Os itens do pedido com o que quem separa precisa: onde achar, quanto tem e o
// nome. Item obrigatorio primeiro, porque e o que nao pode faltar.
func SolicitacaoItens(ctx context.Context, id int64) ([]SolicitacaoItem, error) {
	conn, err := Connection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `
		select
			i.produto,
			coalesce(pr.nome, i.produto) as nome,
			i.quantidade,
			-- Nulo enquanto o item nao foi separado. Sem o coalesce o Scan
			-- quebra, porque o destino e um int e nao um NullInt64.
			coalesce(i.quantidade_separada, 0),
			i.obrigatorio,
			i.separado,
			coalesce(l.nome, '') as local,
			coalesce(pr.endereco, '') as endereco,
			coalesce(pr.quantidade, 0) as saldo
		from solicitacao_item i
		left join epiproduto pr on pr.id = i.produto
		left join local l       on l.id = pr.local
		where i.solicitacao = $1
		order by i.obrigatorio desc, nome
	`, id)
	if err != nil {
		return nil, err
	}

	list := []SolicitacaoItem{}
	defer rows.Close()
	for rows.Next() {
		item := SolicitacaoItem{}

		err := rows.Scan(
			&item.Produto, &item.Nome, &item.Quantidade, &item.QuantidadeSeparada,
			&item.Obrigatorio, &item.Separado, &item.Local, &item.Endereco, &item.Saldo,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

// Marca ou desmarca um item no picking e acerta o status do pedido na mesma
// transacao.
//
// O status e consequencia do que foi separado, nao um campo que a tela controla:
// marcar o primeiro item poe o pedido em separando, marcar o ultimo poe em
// aguardando retirada, e desmarcar volta atras. Assim a cor da linha na Lista de
// Separacao nunca fica mentindo sobre o estado real do pedido.
// quantidade e quanto saiu de fato. Zero ou negativo quando esta desmarcando,
// e nesse caso o campo volta a zero.
func SetItemSeparado(ctx context.Context, id int64, produto string, separado bool, quantidade int, separador string) error {
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

	if !separado || quantidade < 0 {
		quantidade = 0
	}

	_, err = tx.ExecContext(ctx, `
		update solicitacao_item
		set
			separado = $3,
			-- Sem quantidade informada, assume que saiu tudo que foi pedido. E o
			-- caso comum: quem separa so digita numero quando falta.
			quantidade_separada = case
				when not $3 then 0
				when $4 > 0 then least($4, solicitacao_item.quantidade)
				else solicitacao_item.quantidade
			end
		where
			solicitacao_item.solicitacao = $1
			and solicitacao_item.produto = $2
	`, id, produto, separado, quantidade)
	if err != nil {
		return err
	}

	total := 0
	separados := 0
	err = tx.QueryRowContext(ctx, `
		select
			count(*),
			count(*) filter (where separado)
		from solicitacao_item
		where solicitacao = $1
	`, id).Scan(&total, &separados)
	if err != nil {
		return err
	}

	status := StatusSeparando
	if separados == 0 {
		status = StatusAberta
	} else if total > 0 && separados == total {
		status = StatusAguardandoRetirada
	}

	// A data da separacao e calculada aqui, e nao no SQL: comparar dois
	// parametros entre si (o status novo com a constante) deixa o Postgres sem
	// saber o tipo e a query e recusada.
	dataSeparacao := sql.NullTime{}
	if status == StatusAguardandoRetirada {
		dataSeparacao = sql.NullTime{Time: time.Now(), Valid: true}
	}

	// Quem pegou o pedido primeiro fica como separador. Nao sobrescreve a cada
	// item para nao trocar de dono se duas pessoas mexerem no mesmo pedido.
	_, err = tx.ExecContext(ctx, `
		update solicitacao
		set
			status = $2,
			separador = case
				when $3 = '' then separador
				when separador is null or separador = '' then $3
				else separador
			end,
			data_separacao = $4
		where solicitacao.id = $1
	`, id, status, separador, dataSeparacao)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Cancelar so vale enquanto ninguem comecou a separar, e quem cancela e o
// solicitante. Depois que o almoxarifado tocou no pedido, a saida dele da lista
// e pelo prazo de retirada, nao por cancelamento.
//
// A trava fica aqui e nao na tela porque a regra e do negocio: qualquer tela
// que chame isto, hoje ou depois, respeita.
func CancelSolicitacao(ctx context.Context, id int64, pessoa string) error {
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

	result, err := tx.ExecContext(ctx, `
		update solicitacao
		set status = $3
		where
			solicitacao.id = $1
			and solicitacao.pessoa = $2
			and solicitacao.status = $4
	`, id, pessoa, StatusCancelada, StatusAberta)
	if err != nil {
		return err
	}

	// Zero linha quer dizer que o pedido nao e da pessoa ou a separacao ja
	// comecou. Nos dois casos o cancelamento nao vale, e quem chamou precisa
	// saber disso em vez de achar que deu certo.
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return ErrCancelamentoNaoPermitido
	}

	return tx.Commit()
}

func SetSolicitacaoStatus(ctx context.Context, id int64, status string) error {
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
		update solicitacao
		set status = $2
		where solicitacao.id = $1
	`, id, status)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Entrega recusada: ou a requisicao nao esta pronta, ou quem confirmou nao e o
// requisitante.
var ErrEntregaNaoPermitida = errors.New("entrega nao permitida")

// Efetiva a entrega da requisicao.
//
// Quem confirma e o proprio requisitante, olhando os itens na tela e se
// identificando pelo reconhecimento facial. Por isso a funcao recebe quem foi
// reconhecido e compara com quem pediu: a conferencia de identidade acontece
// aqui, nao so na tela.
//
// E aqui que o prazo de troca do EPI zera. Ate agora a retirada nao era
// registrada em lugar nenhum do fluxo de almoxarifado, e o item so parava de ser
// cobrado por existir requisicao em aberto. Com a entrega, o que foi recebido
// entra no historico e a contagem recomeca do zero.
func EntregarSolicitacao(ctx context.Context, id int64, pessoaConfirmada string) error {
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

	// Trava a linha: duas confirmacoes ao mesmo tempo entregariam duas vezes e
	// gravariam a retirada em dobro.
	pessoa := ""
	status := ""
	err = tx.QueryRowContext(ctx, `
		select pessoa, status
		from solicitacao
		where id = $1
		for update
	`, id).Scan(&pessoa, &status)
	if err != nil {
		return err
	}

	if status != StatusAguardandoRetirada || pessoa != pessoaConfirmada {
		return ErrEntregaNaoPermitida
	}

	// Uma retirada por item separado, com a quantidade que saiu de fato. Item
	// que o almoxarifado nao conseguiu separar nao entra: nao foi entregue,
	// entao nao zera prazo nenhum.
	_, err = tx.ExecContext(ctx, `
		INSERT INTO retirada (pessoa, produto, origem)
		SELECT $2, produto, 'almoxarifado'
		FROM solicitacao_item
		WHERE
			solicitacao = $1
			and separado
			and coalesce(quantidade_separada, 0) > 0
	`, id, pessoa)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		update solicitacao
		set status = $2
		where id = $1
	`, id, StatusEntregue)
	if err != nil {
		return err
	}

	return tx.Commit()
}
