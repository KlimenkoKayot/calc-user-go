package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Load загружает конфигурацию, находя корень проекта по .projectroot
func Load() (*Config, error) {
	// 1. Находим корень проекта
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %w", err)
	}

	// 2. Формируем путь к конфигу
	configPath := filepath.Join(projectRoot, "config", "config.yaml")

	// 3. Настраиваем Viper
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 4. Читаем конфиг
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w (tried path: %s)", err, configPath)
	}

	// 5. Подключаем переменные окружения
	v.AutomaticEnv()

	// 6. Парсим в структуру
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// findProjectRoot ищет вверх по директориям до нахождения .projectroot
func findProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		// Проверяем наличие маркера
		markerPath := filepath.Join(currentDir, ".projectroot")
		if _, err := os.Stat(markerPath); err == nil {
			return currentDir, nil
		}

		// Поднимаемся на уровень выше
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Достигли корня файловой системы
			return "", fmt.Errorf(".projectroot not found (reached filesystem root)")
		}
		currentDir = parentDir
	}
}
