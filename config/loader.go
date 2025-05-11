package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Load(path string) (*Config, error) {
	v := viper.New()

	// Устанавливаем значения по умолчанию
	v.SetDefault("common.environment", "development")
	v.SetDefault("common.log_level", "debug")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	v.AddConfigPath(".")

	// Чтение конфигурационного файла
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// Поддержка переменных окружения
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}
	return filepath.Clean(configPath)
}
