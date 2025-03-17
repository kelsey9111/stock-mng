package initialize

import (
	"stock-management/global"
	"strconv"

	"go.uber.org/zap"
)

func Run() {
	defer global.Logger.Sync()
	InitLoadConfig()
	InitLogger()
	global.Logger.Info("Logger initialized", zap.String("status", "success"))
	InitPostgreSQL()
	InitRedis()
	r := InitRouter()
	port := strconv.Itoa(global.Config.Server.Port)
	if err := r.Run(":" + port); err != nil {
		global.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
