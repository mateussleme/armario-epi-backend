package door

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/TopSisErp/epi-backend/internal/viaonda"
	"github.com/gin-gonic/gin"
)

func Routes(group *gin.RouterGroup) {
	door := group.Group("/door")

	door.POST("/open", func(c *gin.Context) {
		portsStr := strings.Split(os.Getenv("DOOR_PORTS"), ",")
		portStates := strings.Split(os.Getenv("DOOR_STATES"), ",")
		ports := make([]int, len(portsStr))

		for i, portStr := range portsStr {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				slog.Error("error converting relay port", "err", err.Error(), "port", portStr, "index", i)
				return
			}

			ports[i] = port
		}

		for index, port := range ports {
			err := viaonda.SetGpio(c, port, portStates[index] != "1")
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				slog.Error("error opening door", "err", err.Error())
				return
			}
		}

		c.Status(http.StatusOK)
	})

	door.POST("/close", func(c *gin.Context) {
		portsStr := strings.Split(os.Getenv("DOOR_PORTS"), ",")
		portStates := strings.Split(os.Getenv("DOOR_STATES"), ",")
		ports := make([]int, len(portsStr))

		for i, portStr := range portsStr {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				slog.Error("error converting relay port", "err", err.Error(), "port", portStr, "index", i)
				return
			}

			ports[i] = port
		}

		for index, port := range ports {
			err := viaonda.SetGpio(c, port, portStates[index] == "1")
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				slog.Error("error closing door", "err", err.Error())
				return
			}
		}

		c.Status(http.StatusOK)
	})

	door.GET("/state", func(c *gin.Context) {
		portsStr := strings.Split(os.Getenv("SENSOR_PORTS"), ",")
		portStates := strings.Split(os.Getenv("SENSOR_STATES"), ",")
		ports := make([]int, len(portsStr))

		for i, portStr := range portsStr {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				slog.Error("error converting sensor port", "err", err.Error(), "port", portStr, "index", i)
				return
			}

			ports[i] = port
		}

		totalState := true
		for index, port := range ports {
			rawState, err := viaonda.GetGpio(c, port)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				slog.Error("error getting door", "err", err.Error())
				return
			}

			state := rawState
			if portStates[index] != "1" {
				state = !state
			}

			totalState = totalState && state
		}

		stateText := "opened"
		if totalState {
			stateText = "closed"
		}

		c.JSON(200, gin.H{
			"state": stateText,
		})
	})
}
