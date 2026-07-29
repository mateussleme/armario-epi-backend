package session

import (
	"log/slog"
	"net/http"

	"github.com/TopSisErp/epi-backend/internal/database"
	"github.com/TopSisErp/epi-backend/internal/viaonda"
	"github.com/gin-gonic/gin"
)

func Routes(group *gin.RouterGroup) {
	session := group.Group("/session")
	session.POST("/initiate", func(c *gin.Context) {
		err := viaonda.ClearEntries(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error clearing data", "err", err.Error())
			return
		}

		c.Status(http.StatusOK)
	})

	session.GET("/poke", func(c *gin.Context) {
		tags, err := viaonda.GetEntries(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error getting data", "err", err.Error())
			return
		}

		products, err := database.TagsToProducts(c, tags)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error translating data", "err", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"products": products,
		})
	})

	session.POST("/commit", func(c *gin.Context) {
		tags, err := viaonda.GetEntries(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error getting data", "err", err.Error())
			return
		}

		err = database.UpdateWithTags(c, tags)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error updating data", "err", err.Error())
			return
		}

		c.Status(http.StatusOK)
	})

	session.GET("/inventory/unknown", func(c *gin.Context) {
		tags, err := viaonda.GetEntries(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error getting data", "err", err.Error())
			return
		}

		newTags, err := database.UnknownTags(c, tags)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error getting unknown tags", "err", err.Error())
			return
		}

		idList := []string{}
		for _, tag := range newTags {
			idList = append(idList, tag.Id)
		}

		c.JSON(http.StatusOK, gin.H{
			"tags": idList,
		})
	})

	session.POST("/inventory/:id", func(c *gin.Context) {
		tags, err := viaonda.GetEntries(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error getting data", "err", err.Error())
			return
		}

		newTags, err := database.UnknownTags(c, tags)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error getting unknown tags", "err", err.Error())
			return
		}

		id := c.Param("id")
		err = database.InventoryWithTags(c, newTags, id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error inventoring data", "err", err.Error())
			return
		}

		c.Status(http.StatusOK)
	})
}
