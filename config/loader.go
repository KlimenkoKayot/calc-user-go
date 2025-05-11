package config

import (
	"bytes"
	"embed"

	"github.com/spf13/viper"
)

//go:embed config.yaml
var configFS embed.FS

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	// Читаем из встроенного файла
	file, err := configFS.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	if err := v.ReadConfig(bytes.NewReader(file)); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
