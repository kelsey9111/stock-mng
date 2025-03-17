package initialize

import (
	"stock-management/global"
	"stock-management/internal/middlewares"
	"stock-management/internal/routers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}
	r.Use(middlewares.RecoveryMiddleware())
	r.Use(middlewares.RequestLogger())

	managerRouter := routers.RouterGroupApp.Manager
	mainGroup := r.Group("api")
	managerRouter.InitAdminRouter(mainGroup)
	return r
}
