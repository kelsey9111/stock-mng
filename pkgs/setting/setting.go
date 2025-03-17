package setting

type SettingConfig struct {
	MySql      MySqlSetting      `mapstructure:"mysql"`
	PostgreSql PostgreSqlSetting `mapstructure:"postgresql"`
	Logger     LoggerSetting     `mapstructure:"logger"`
	Server     ServerSetting     `mapstructure:"server"`
	Cache      RedisSetting      `mapstructure:"redis"`
}

type RedisSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DBname   int    `mapstructure:"dbname"`
	PoolSize int    `mapstructure:"poolSize"`
}

type MySqlSetting struct {
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	Username         string `mapstructure:"username"`
	Password         string `mapstructure:"password"`
	DBname           string `mapstructure:"dbname"`
	MaxIdleConns     int    `mapstructure:"maxIdleConns"`
	MaxOpenConns     int    `mapstructure:"maxOpenConns"`
	MaxLifeTimeConns int    `mapstructure:"maxLifeTimeConns"`
}

type PostgreSqlSetting struct {
	Host             string `mapstructure:"host"`
	Port             int    `mapstructure:"port"`
	Username         string `mapstructure:"username"`
	Password         string `mapstructure:"password"`
	DBname           string `mapstructure:"dbname"`
	MaxIdleConns     int    `mapstructure:"maxIdleConns"`
	MaxOpenConns     int    `mapstructure:"maxOpenConns"`
	MaxLifeTimeConns int    `mapstructure:"maxLifeTimeConns"`
}

type LoggerSetting struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"fileName"`
	MaxSize    int    `mapstructure:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge"`
	Compress   bool   `mapstructure:"compress"`
}

type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}
