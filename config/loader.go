package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

//go:embed config.yaml
var configFS embed.FS

func findProjectRoot() (string, error) {
	current, _ := os.Getwd()

	for {
		if _, err := os.Stat(filepath.Join(current, ".projectroot")); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("project root not found")
		}
		current = parent
	}
}

func Load() (*Config, error) {
	root, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigFile(filepath.Join(root, "config", "config.yaml"))

	fmt.Println(filepath.Join(root, "config", "config.yaml"))

	// Устанавливаем значения по умолчанию
	v.SetDefault("common.environment", "development")
	v.SetDefault("common.log_level", "debug")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
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
