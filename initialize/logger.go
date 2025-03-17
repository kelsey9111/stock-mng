package initialize

import (
	"stock-management/global"
	"stock-management/pkgs/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
