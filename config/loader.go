package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Load загружает конфигурацию из YAML файла с учетом правильного пути
func Load() (*Config, error) {
	v := viper.New()

	// 1. Определяем абсолютный путь к config.yaml
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// 2. Настраиваем Viper
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 3. Читаем конфиг
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file at %s: %w", configPath, err)
	}

	// 4. Подключаем переменные окружения
	v.AutomaticEnv()

	// 5. Выводим информацию о загруженном конфиге (для отладки)
	log.Printf("Successfully loaded config from: %s", v.ConfigFileUsed())

	// 6. Парсим конфиг в структуру
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// getConfigPath возвращает абсолютный путь к config.yaml
func getConfigPath() (string, error) {
	// Вариант 1: Используем переменную окружения
	if customPath := os.Getenv("CONFIG_PATH"); customPath != "" {
		if filepath.IsAbs(customPath) {
			return customPath, nil
		}
		return filepath.Abs(customPath)
	}

	// Вариант 2: Автоматический поиск относительно исполняемого файла
	_, currentFilePath, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFilePath)))
	configPath := filepath.Join(projectRoot, "config", "config.yaml")

	// Проверяем существование файла
	if _, err := os.Stat(configPath); err == nil {
		return configPath, nil
	}

	// Вариант 3: Поиск в текущей директории
	currentDir, _ := os.Getwd()
	localConfigPath := filepath.Join(currentDir, "config.yaml")
	if _, err := os.Stat(localConfigPath); err == nil {
		return localConfigPath, nil
	}

	return "", fmt.Errorf("config file not found in: %s or %s", configPath, localConfigPath)
}
