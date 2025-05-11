package config

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

//go:embed config.yaml
var embeddedConfig embed.FS

// findProjectRoot ищет корневую директорию проекта по маркерному файлу
func findProjectRoot(marker string) (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	for {
		// Проверяем наличие маркерного файла
		markerPath := filepath.Join(current, marker)
		if _, err := os.Stat(markerPath); err == nil {
			return current, nil
		}

		// Поднимаемся на уровень выше
		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("project root not found (reached filesystem root)")
		}
		current = parent
	}
}

// Load загружает конфигурацию, автоматически находя корень проекта
func Load() (*Config, error) {
	// 1. Пытаемся найти корень проекта
	root, err := findProjectRoot(".projectroot")
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %w", err)
	}

	// 2. Пробуем загрузить из файла
	cfgPath := filepath.Join(root, "config", "config.yaml")
	v := viper.New()
	v.SetConfigFile(cfgPath)

	if err := v.ReadInConfig(); err == nil {
		// Успешно загрузили из файла
		return unmarshalConfig(v)
	}

	// 3. Если файл не найден, пробуем встроенную конфигурацию
	file, err := embeddedConfig.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded config: %w", err)
	}

	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader(file)); err != nil {
		return nil, fmt.Errorf("failed to read embedded config: %w", err)
	}

	return unmarshalConfig(v)
}

// unmarshalConfig выполняет общие шаги для unmarshal конфига
func unmarshalConfig(v *viper.Viper) (*Config, error) {
	// Поддержка переменных окружения
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
