package initialize

import (
	"fmt"
	"stock-management/global"

	"github.com/spf13/viper"
)

func InitLoadConfig() {
	viper := viper.New()
	viper.AddConfigPath("./config/")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fail to read config %v", err))
	}
	if err := viper.Unmarshal(&global.Config); err != nil {
		fmt.Printf("unable decode config %v \n", err)
	}
}
