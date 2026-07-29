package item

import (
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/TopSisErp/epi-backend/internal/database"
	"github.com/gin-gonic/gin"
)

type CreateForm struct {
	Name        string                `form:"name"`
	Description string                `form:"description"`
	Image       *multipart.FileHeader `form:"image"`
	VideoUri    string                `form:"videoUri"`
}

func Routes(group *gin.RouterGroup) {
	items := group.Group("/items")

	items.GET("/:id/data", func(c *gin.Context) {
		id := c.Param("id")

		item, err := database.ItemData(c, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}

			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error getting item", "err", err.Error())
			return
		}

		c.JSON(200, gin.H{
			"item": item,
		})
	})

	items.GET("/all", func(c *gin.Context) {
		data, err := database.AllItems(c)
		if err != nil {
			slog.Error("error getting items", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"items": data,
		})
	})

	items.POST("/:id/create", func(c *gin.Context) {
		id := c.Param("id")

		var newItem CreateForm
		if err := c.Bind(&newItem); err != nil {
			slog.Error("error binding", "err", err.Error())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var imageBytes []byte
		if newItem.Image != nil {
			file, err := newItem.Image.Open()
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

		err := database.UpdateItem(c, id, newItem.Name, newItem.Description, imageBytes, newItem.VideoUri)
		if err != nil {
			slog.Error("error updating item", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})

	items.POST("/:id/delete", func(c *gin.Context) {
		id := c.Param("id")

		err := database.DeleteItem(c, id)
		if err != nil {
			slog.Error("error deleting item", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	})
}
