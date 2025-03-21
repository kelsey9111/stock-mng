package initialize

import (
	_ "stock-management/docs"
	"stock-management/global"
	"strconv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func Run() {
	defer global.Logger.Sync()
	InitLoadConfig()
	InitLogger()
	global.Logger.Info("Logger initialized", zap.String("status", "success"))
	InitPostgreSQL()
	InitRedis()
	InitService()
	r := InitRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	port := strconv.Itoa(global.Config.Server.Port)
	if err := r.Run(":" + port); err != nil {
		global.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
