package global

import (
	"stock-management/pkgs/logger"
	"stock-management/pkgs/setting"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config setting.SettingConfig
	Logger *logger.LoggerZap
	Pdb    *gorm.DB
	Rdb    *redis.Client
)
