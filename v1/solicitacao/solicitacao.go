package solicitacao

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/TopSisErp/epi-backend/internal/database"
	"github.com/gin-gonic/gin"
)

// Diferente do resto do projeto, esta rota recebe JSON e nao multipart. Um
// pedido e um cabecalho com uma lista de itens dentro, e lista em form viraria
// aquele itens[0][produto] que ninguem consegue ler depois.
type CreateForm struct {
	Pessoa string `json:"pessoa"`
	Itens  []struct {
		Produto     string `json:"produto"`
		Quantidade  int    `json:"quantidade"`
		Obrigatorio bool   `json:"obrigatorio"`
	} `json:"itens"`
}

func parseId(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		slog.Error("error parsing solicitacao id", "err", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func Routes(group *gin.RouterGroup) {
	solicitacoes := group.Group("/solicitacoes")

	solicitacoes.POST("/create", func(c *gin.Context) {
		var form CreateForm
		if err := c.ShouldBindJSON(&form); err != nil {
			slog.Error("error binding", "err", err.Error())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if form.Pessoa == "" || len(form.Itens) == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		itens := []database.NovaSolicitacaoItem{}
		for _, item := range form.Itens {
			if item.Produto == "" {
				continue
			}
			itens = append(itens, database.NovaSolicitacaoItem{
				Produto:     item.Produto,
				Quantidade:  item.Quantidade,
				Obrigatorio: item.Obrigatorio,
			})
		}

		if len(itens) == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		id, err := database.CreateSolicitacao(c, form.Pessoa, itens)
		if err != nil {
			slog.Error("error creating solicitacao", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	// A Lista de Separacao. Sem status na query traz tudo que ainda esta em
	// aberto, que e como a tela abre.
	solicitacoes.GET("/all", func(c *gin.Context) {
		list, err := database.AllSolicitacoes(c, c.Query("status"))
		if err != nil {
			slog.Error("error getting solicitacoes", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"solicitacoes": list})
	})

	solicitacoes.GET("/:id/data", func(c *gin.Context) {
		id, ok := parseId(c)
		if !ok {
			return
		}

		data, err := database.SolicitacaoData(c, id)
		if err != nil {
			slog.Error("error getting solicitacao", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"solicitacao": data})
	})

	// Os itens para a tela de picking: nome, local, endereco e saldo.
	solicitacoes.GET("/:id/itens", func(c *gin.Context) {
		id, ok := parseId(c)
		if !ok {
			return
		}

		itens, err := database.SolicitacaoItens(c, id)
		if err != nil {
			slog.Error("error getting solicitacao itens", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{"itens": itens})
	})

	// Marca ou desmarca um item. O status do pedido sai daqui, calculado pelo
	// banco, entao a tela nao precisa mandar status nenhum.
	solicitacoes.PUT("/:id/item/:produto", func(c *gin.Context) {
		id, ok := parseId(c)
		if !ok {
			return
		}

		separado := c.Query("separado") != "0"
		// Sem quantidade na query, o banco assume que saiu tudo que foi pedido.
		quantidade, _ := strconv.Atoi(c.Query("quantidade"))

		err := database.SetItemSeparado(c, id, c.Param("produto"), separado, quantidade, c.Query("separador"))
		if err != nil {
			slog.Error("error setting item separado", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})

	// Cancelamento pelo solicitante. Pede a pessoa na query porque o cancelamento
	// e dela: o backend confere que o pedido e mesmo daquele usuario e que a
	// separacao ainda nao comecou, em vez de confiar na tela.
	//
	// O pedido nao e apagado. Ele some da lista, mas o registro de que alguem
	// pediu e desistiu continua.
	solicitacoes.POST("/:id/cancel", func(c *gin.Context) {
		id, ok := parseId(c)
		if !ok {
			return
		}

		pessoa := c.Query("pessoa")
		if pessoa == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err := database.CancelSolicitacao(c, id, pessoa)
		if errors.Is(err, database.ErrCancelamentoNaoPermitido) {
			// Nao e erro de servidor: a regra recusou.
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		if err != nil {
			slog.Error("error canceling solicitacao", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})
}
