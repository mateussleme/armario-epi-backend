package v1

import (
	"github.com/TopSisErp/epi-backend/v1/door"
	"github.com/TopSisErp/epi-backend/v1/groups"
	"github.com/TopSisErp/epi-backend/v1/item"
	"github.com/TopSisErp/epi-backend/v1/local"
	"github.com/TopSisErp/epi-backend/v1/session"
	"github.com/TopSisErp/epi-backend/v1/solicitacao"
	"github.com/TopSisErp/epi-backend/v1/stock"
	"github.com/TopSisErp/epi-backend/v1/user"
	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine) {
	v1 := engine.Group("/v1")

	user.Routes(v1)
	item.Routes(v1)
	groups.Routes(v1)
	local.Routes(v1)
	door.Routes(v1)
	session.Routes(v1)
	solicitacao.Routes(v1)
	stock.Routes(v1)
}
