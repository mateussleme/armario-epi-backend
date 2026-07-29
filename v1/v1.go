package v1

import (
	"log/slog"
	"net/http"

	"github.com/TopSisErp/epi-backend/internal/database"
	"github.com/TopSisErp/epi-backend/v1/door"
	"github.com/TopSisErp/epi-backend/v1/groups"
	"github.com/TopSisErp/epi-backend/v1/item"
	"github.com/TopSisErp/epi-backend/v1/session"
	"github.com/TopSisErp/epi-backend/v1/user"
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {
	v1 := engine.Group("/v1")

	v1.GET("/stock", func(c *gin.Context) {
		stock, err := database.GetStock(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			slog.Error("error clearing data", "err", err.Error())
			return
		}

		c.JSON(200, gin.H{
			"products": stock,
		})
	})

	door.Routes(v1)
	user.Routes(v1)
	groups.Routes(v1)
	item.Routes(v1)
	session.Routes(v1)
}
