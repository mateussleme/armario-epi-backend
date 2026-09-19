package database

type Item struct {
	Id           string `json:"id"`
	Name         string `json:"title"`
	Description  string `json:"description"`
	ImageUri     string `json:"imageUri"`
	VideoUri     string `json:"videoUri"`
	Ca           string `json:"ca"`
	CaVencimento string `json:"caVencimento"`
	Endereco     string `json:"endereco"`
	Porta        int    `json:"porta"`
	Local        string `json:"local"`
}

type Local struct {
	Id      string `json:"id"`
	Nome    string `json:"nome"`
	Tipo    string `json:"tipo"`
	Empresa string `json:"empresa"`
	Ativo   bool   `json:"ativo"`
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
	ImageUri string `json:"imageUri"`
}

type Group struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// Vinculo entre grupo e produto. O prazo mora aqui, e nao no produto, porque o
// desgaste depende da atividade do GES: a mesma luva dura 5 dias na eletrica e
// pode durar 15 na logistica.
type GroupProduct struct {
	Produto      string `json:"produto"`
	Obrigatorio  bool   `json:"obrigatorio"`
	// Nulo quando o item nao vence. Ponteiro para distinguir "sem prazo" de
	// "prazo zero", que sao coisas diferentes na regra.
	DiasValidade *int   `json:"diasValidade"`
}

// Item que a pessoa e obrigada a pedir: obrigatorio no GES e vencido (ou nunca
// retirado). Vai para o carrinho travado.
type RequiredProduct struct {
	Produto       string `json:"produto"`
	DiasValidade  int    `json:"diasValidade"`
	// Dias desde a ultima retirada. -1 quando nunca retirou.
	DiasDesdeUso  int    `json:"diasDesdeUso"`
	UltimaRetirada string `json:"ultimaRetirada"`
}

type Retirada struct {
	Pessoa      string `json:"pessoa"`
	Produto     string `json:"produto"`
	// Nome resolvido no momento da consulta. Vazio quando o produto ja foi
	// excluido do cadastro: o historico continua valendo, so perde o nome.
	ProdutoNome string `json:"produtoNome"`
	Quantidade  int    `json:"quantidade"`
	Data        string `json:"data"`
	Origem      string `json:"origem"`
}

// Uma linha da Lista de Separacao. Ja traz o nome de quem pediu e a contagem de
// itens porque e exatamente o que o grid mostra: sem isto a tela faria uma
// consulta por linha so para escrever o nome.
type Solicitacao struct {
	Id            int64  `json:"id"`
	Pessoa        string `json:"pessoa"`
	PessoaNome    string `json:"pessoaNome"`
	Data          string `json:"data"`
	// Hora limite para separar. A tela pinta a linha a partir disto.
	SepararAte    string `json:"separarAte"`
	Status        string `json:"status"`
	Separador     string `json:"separador"`
	DataSeparacao string `json:"dataSeparacao"`
	TotalItens    int    `json:"totalItens"`
	ItensSeparados int   `json:"itensSeparados"`
}

// Item dentro do pedido, na tela de picking. Traz local, saldo e foto porque e
// o que o Clairton pediu para quem esta separando: onde achar, quanto tem e
// como o item e.
type SolicitacaoItem struct {
	Produto     string `json:"produto"`
	Nome        string `json:"nome"`
	Quantidade  int    `json:"quantidade"`
	// Quanto saiu de fato. Menor que Quantidade quando faltou saldo no local.
	QuantidadeSeparada int `json:"quantidadeSeparada"`
	Obrigatorio bool   `json:"obrigatorio"`
	Separado    bool   `json:"separado"`
	Local       string `json:"local"`
	Endereco    string `json:"endereco"`
	Saldo       int    `json:"saldo"`
}

// A foto nao vem aqui de proposito. Ela e bytea no banco e viraria um data URL
// grande por item; a tela busca por item com o mesmo componente que as outras
// telas ja usam.

const (
	StatusAberta             = "aberta"
	StatusSeparando          = "separando"
	StatusAguardandoRetirada = "aguardando_retirada"
	StatusCancelada          = "cancelada"
	StatusEntregue           = "entregue"
)
