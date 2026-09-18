package local

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/TopSisErp/epi-backend/internal/database"
	"github.com/gin-gonic/gin"
)

type CreateForm struct {
	Nome  string `form:"nome"`
	Tipo  string `form:"tipo"`
	Ativo string `form:"ativo"`
}

func Routes(group *gin.RouterGroup) {
	locais := group.Group("/locais")

	locais.GET("/:id/data", func(c *gin.Context) {
		id := c.Param("id")

		local, err := database.LocalData(c, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			slog.Error("error getting local", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"local": local,
		})
	})

	locais.GET("/all", func(c *gin.Context) {
		data, err := database.AllLocais(c)
		if err != nil {
			slog.Error("error getting locais", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"locais": data,
		})
	})

	locais.POST("/:id/create", func(c *gin.Context) {
		id := c.Param("id")

		var newLocal CreateForm
		if err := c.Bind(&newLocal); err != nil {
			slog.Error("error binding", "err", err.Error())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// O formulario manda "1" ou "0"; qualquer coisa diferente de "0" conta
		// como ativo, para um local novo nascer disponivel.
		ativo := newLocal.Ativo != "0"

		err := database.UpdateLocal(c, id, newLocal.Nome, newLocal.Tipo, ativo)
		if err != nil {
			slog.Error("error updating local", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	locais.POST("/:id/delete", func(c *gin.Context) {
		id := c.Param("id")

		err := database.DeleteLocal(c, id)
		if err != nil {
			// O caso comum aqui e ter produto apontando para o local. O banco
			// recusa por causa da chave estrangeira, e a tela orienta a
			// desativar em vez de excluir.
			slog.Error("error deleting local", "err", err.Error())
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})
}
