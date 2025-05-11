package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   string
	Router   string
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func Load(path string) (*Config, error) {
	viper.SetConfigName("config") // имя файла конфигурации без расширения
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".") // дополнительно ищем в текущей директории

	// Устанавливаем значения по умолчанию для SQLite
	viper.SetDefault("database.dsn", "file:auth.db?cache=shared&mode=rwc")

	// Чтение конфигурации из файла
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Поддержка переменных окружения (переопределяют значения из файла)
	viper.AutomaticEnv()

	// Для совместимости можно оставить поддержку .env файлов
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.MergeInConfig(); err != nil {
			return nil, err
		}
	}

	config := &Config{
		Logger: viper.GetString("logger"),
		Router: viper.GetString("router"),
		Server: ServerConfig{
			Port:         viper.GetInt("server.port"),
			ReadTimeout:  viper.GetDuration("server.read_timeout"),
			WriteTimeout: viper.GetDuration("server.write_timeout"),
		},
		Database: DatabaseConfig{
			DSN: viper.GetString("database.dsn"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("jwt.secret"),
			AccessTokenExpiry:  viper.GetDuration("jwt.access_token_expiration"),
			RefreshTokenExpiry: viper.GetDuration("jwt.refresh_token_expiration"),
		},
	}

	return config, nil
}
