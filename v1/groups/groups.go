package groups

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/TopSisErp/epi-backend/internal/database"
	"github.com/gin-gonic/gin"
)

type CreateForm struct {
	Name string `form:"name"`
}

// Regra do item dentro do grupo. Vem como texto porque o front manda multipart,
// igual ao resto das telas de cadastro.
//
// DiasValidade vazio significa "nao vence", que e diferente de zero.
type ItemForm struct {
	Obrigatorio  string `form:"obrigatorio"`
	DiasValidade string `form:"diasValidade"`
}

func Routes(group *gin.RouterGroup) {
	groups := group.Group("/groups")

	groups.GET("/:id/data", func(c *gin.Context) {
		id := c.Param("id")

		data, err := database.GroupData(c, id)
		if err != nil {
			slog.Error("error getting group", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"group": data,
		})
	})

	groups.GET("/all", func(c *gin.Context) {
		data, err := database.AllGroups(c)
		if err != nil {
			slog.Error("error getting groups", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"groups": data,
		})
	})

	// Os produtos do grupo com a regra de cada um (obrigatorio, prazo de troca).
	groups.GET("/:id/products", func(c *gin.Context) {
		id := c.Param("id")

		products, err := database.GroupProducts(c, id)
		if err != nil {
			slog.Error("error getting group products", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"products": products,
		})
	})

	groups.POST("/:id/create", func(c *gin.Context) {
		id := c.Param("id")

		var newGroup CreateForm
		if err := c.Bind(&newGroup); err != nil {
			slog.Error("error binding", "err", err.Error())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err := database.UpdateGroup(c, id, newGroup.Name)
		if err != nil {
			slog.Error("error updating group", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	groups.POST("/:id/delete", func(c *gin.Context) {
		id := c.Param("id")

		err := database.DeleteGroup(c, id)
		if err != nil {
			slog.Error("error deleting group", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	// Adiciona o item ao grupo e grava a regra. Serve tambem para so editar a
	// regra de um item que ja esta la: o banco faz INSERT ... ON CONFLICT UPDATE.
	groups.PUT("/:id/item/:product", func(c *gin.Context) {
		id := c.Param("id")
		product := c.Param("product")

		var form ItemForm
		// Bind sem erro fatal: quem so quer vincular o item, sem regra, manda a
		// requisicao vazia e cai nos defaults.
		_ = c.ShouldBind(&form)

		obrigatorio := form.Obrigatorio == "1" || form.Obrigatorio == "true"

		var dias *int
		if form.DiasValidade != "" {
			value, err := strconv.Atoi(form.DiasValidade)
			if err != nil {
				slog.Error("error parsing diasValidade", "err", err.Error())
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			if value > 0 {
				dias = &value
			}
		}

		err := database.GroupSetProduct(c, id, product, obrigatorio, dias)
		if err != nil {
			slog.Error("error adding to group", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})

	groups.DELETE("/:id/item/:product", func(c *gin.Context) {
		id := c.Param("id")
		product := c.Param("product")

		err := database.GroupRemoveProduct(c, id, product)
		if err != nil {
			slog.Error("error removing from group", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})
}
