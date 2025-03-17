package routers

import (
	"stock-management/internal/routers/manager"
)

type RouterGroup struct {
	Manager manager.ManageRouterGroup
}

var RouterGroupApp = new(RouterGroup)
