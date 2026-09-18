package user

import (
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/TopSisErp/epi-backend/internal/database"
	"github.com/gin-gonic/gin"
)

type CreateForm struct {
	Name  string                `form:"name"`
	Admin string                `form:"admin"`
	Image *multipart.FileHeader `form:"image"`
}

func Routes(group *gin.RouterGroup) {
	user := group.Group("/user")

	user.POST("/:id/create", func(c *gin.Context) {
		id := c.Param("id")

		var newUser CreateForm
		if err := c.Bind(&newUser); err != nil {
			slog.Error("error binding", "err", err.Error())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var imageBytes []byte
		if newUser.Image != nil {
			file, err := newUser.Image.Open()
			if err != nil {
				slog.Error("error opening file", "err", err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer file.Close()
			imageBytes, err = io.ReadAll(file)
			if err != nil {
				slog.Error("error reading file", "err", err.Error())
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		admin := newUser.Admin == "1"
		err := database.UpdateUser(c, id, newUser.Name, imageBytes, admin)
		if err != nil {
			slog.Error("error updating user", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	user.POST("/:id/delete", func(c *gin.Context) {
		id := c.Param("id")

		err := database.DeleteUser(c, id)
		if err != nil {
			slog.Error("error deleting user", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	user.GET("/:id/data", func(c *gin.Context) {
		id := c.Param("id")

		data, err := database.UserData(c, id)
		if err != nil {
			slog.Error("error getting user", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		products, err := database.UserProducts(c, id)
		if err != nil {
			slog.Error("error getting user products", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// requiredItems existia na resposta mas vinha sempre vazio. Agora traz os
		// itens obrigatorios que venceram (ou que a pessoa nunca retirou): sao os
		// que o front coloca no carrinho travado.
		required, err := database.UserRequiredProducts(c, id)
		if err != nil {
			slog.Error("error getting required products", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		requiredIds := []string{}
		for _, item := range required {
			requiredIds = append(requiredIds, item.Produto)
		}

		c.JSON(http.StatusOK, gin.H{
			"user":           data,
			"requiredItems":  requiredIds,
			"availableItems": products,
			// Detalhe do vencimento (prazo e dias desde o ultimo uso), para a tela
			// explicar por que o item esta travado.
			"requiredDetail": required,
		})
	})

	// Registra que a pessoa levou o item. E o que zera a contagem do prazo.
	user.POST("/:id/retirada/:product", func(c *gin.Context) {
		id := c.Param("id")
		product := c.Param("product")
		origem := c.Query("origem")

		err := database.RegisterRetirada(c, id, product, origem)
		if err != nil {
			slog.Error("error registering retirada", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})

	// Historico de entrega, para auditoria.
	user.GET("/:id/retiradas", func(c *gin.Context) {
		id := c.Param("id")

		limit, _ := strconv.Atoi(c.Query("limit"))
		list, err := database.UserRetiradas(c, id, limit)
		if err != nil {
			slog.Error("error getting retiradas", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"retiradas": list,
		})
	})

	user.GET("/all", func(c *gin.Context) {
		data, err := database.AllUsers(c)
		if err != nil {
			slog.Error("error getting users", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"users": data,
		})
	})

	user.PUT("/:id/group/:group", func(c *gin.Context) {
		id := c.Param("id")
		group := c.Param("group")

		err := database.UserAddGroup(c, id, group)
		if err != nil {
			slog.Error("error adding to group", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})

	user.DELETE("/:id/group/:group", func(c *gin.Context) {
		id := c.Param("id")
		group := c.Param("group")

		err := database.UserRemoveGroup(c, id, group)
		if err != nil {
			slog.Error("error removing from group", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})
}
