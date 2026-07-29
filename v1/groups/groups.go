package groups

import (
	"log/slog"
	"net/http"

	"github.com/TopSisErp/epi-backend/internal/database"
	"github.com/gin-gonic/gin"
)

type CreateForm struct {
	Name string `form:"name"`
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

	groups.PUT("/:id/item/:product", func(c *gin.Context) {
		id := c.Param("id")
		product := c.Param("product")

		err := database.GroupAddProduct(c, id, product)
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
