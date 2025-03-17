package initialize

import (
	"context"
	"fmt"
	"stock-management/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitRedis() {
	m := global.Config.Cache
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", m.Host, m.Port),
		Password: m.Password,
		DB:       m.DBname,
		PoolSize: m.PoolSize,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		global.Logger.Error("Failed to connect to redis", zap.Error(err))
	}
	global.Logger.Info("Init redis success")
	global.Rdb = rdb
}
