package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AuthURL string
	CalcURL string
	Port    int
	Router  string
	Logger  string
}

func Load(path string) (*Config, error) {
	// Устанавливаем значения по умолчанию
	viper.SetDefault("auth_url", "http://localhost:8081")
	viper.SetDefault("calc_url", "http://localhost:8082")
	viper.SetDefault("port", 8080)
	viper.SetDefault("router", "gorilla")
	viper.SetDefault("logger", "zap")

	viper.SetConfigName("config") // имя файла без расширения
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".") // ищем в текущей директории

	// Чтение конфигурационного файла
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	viper.AutomaticEnv()
	viper.BindEnv("auth_url")
	viper.BindEnv("calc_url")
	viper.BindEnv("port")
	viper.BindEnv("router")
	viper.BindEnv("logger")

	return &Config{
		AuthURL: viper.GetString("auth_url"),
		CalcURL: viper.GetString("calc_url"),
		Port:    viper.GetInt("port"),
		Router:  viper.GetString("router"),
		Logger:  viper.GetString("logger"),
	}, nil
}
