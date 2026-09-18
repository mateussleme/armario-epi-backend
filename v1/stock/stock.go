package stock

import (
	"log/slog"
	"net/http"

	"github.com/TopSisErp/epi-backend/internal/database"
	"github.com/gin-gonic/gin"
)

func Routes(group *gin.RouterGroup) {
	stock := group.Group("/stock")

	stock.GET("", func(c *gin.Context) {
		products, err := database.GetStock(c)
		if err != nil {
			slog.Error("error getting stock", "err", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"products": products,
		})
	})
}
