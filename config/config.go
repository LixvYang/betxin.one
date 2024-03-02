package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`
	Driver    string `mapstrucutre:"driver"`

	LogConfig   *LogConfig   `mapstructure:"log"`
	MySQLConfig *MySQLConfig `mapstructure:"mysql"`
	MongoConfig *MongoConfig `mapstructure:"mongo"`
	RedisConfig *RedisConfig `mapstructure:"redis"`
	MixinConfig *MixinConfig `mapstructure:"mixin"`
}

type MixinConfig struct {
	ClientID   string `mapstructure:"client_id"`
	SessionID  string `mapstructure:"session_id"`
	PrivateKey string `mapstructure:"private_key"`
	PinToken   string `mapstructure:"pin_token"`
	Scope      string `mapstructure:"scope"`

	// AppID is equivalent to the ClientID
	AppID string `mapstructure:"app_id"`
	// ServerPublicKey is equivalent to the PinToken in hex format
	ServerPublicKey string `mapstructure:"server_public_key"`
	// SessionPrivateKey is equivalent to the PrivateKey in hex format
	SessionPrivateKey string `mapstructure:"session_private_key"`
	SpendKey          string `mapstructure:"spend_key"`
}

type MongoConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       string `mapstructure:"db"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// Configuration for logging
type LogConfig struct {
	ConsoleLoggingEnabled bool   `mapstructure:"console_logging_enabled"`
	EncodeLogsAsJson      bool   `mapstructure:"encode_logs_as_json"`
	FileLoggingEnabled    bool   `mapstructure:"file_logging_enabled"`
	Directory             string `mapstructure:"directory"`
	Filename              string `mapstructure:"filename"`
	MaxSize               int    `mapstructure:"max_size"`
	MaxBackups            int    `mapstructure:"max_backups"`
	MaxAge                int    `mapstructure:"max_age"`
	Level                 int    `mapstructure:"level"`
	LocalTime             bool   `mapstructure:"local_time"`
	Compress              bool   `mapstructure:"compress"`
}

func Init(filePath string, conf *AppConfig) (err error) {
	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		log.Panic().Msgf("viper.ReadInConfig failed, err:%v\n", err)
		return err
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err = viper.Unmarshal(conf); err != nil {
		log.Panic().Msgf("viper.Unmarshal failed, err:%v\n", err)
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(_ fsnotify.Event) {
		log.Info().Msg("配置文件修改了")
		if err := viper.Unmarshal(conf); err != nil {
			log.Info().Msgf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return nil
}
