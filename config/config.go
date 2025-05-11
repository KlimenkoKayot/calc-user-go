package config

import "time"

type Config struct {
	ApiGateway struct {
		Http struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"http"`
		Services struct {
			Auth struct {
				URL string `mapstructure:"url"`
			} `mapstructure:"auth"`
			Calc struct {
				URL string `mapstructure:"url"`
			} `mapstructure:"calc"`
		} `mapstructure:"services"`
		Router string `mapstructure:"router"`
		Logger string `mapstructure:"logger"`
	} `mapstructure:"api_gateway"`

	Auth struct {
		Http struct {
			Host         string        `mapstructure:"host"`
			Port         int           `mapstructure:"port"`
			ReadTimeout  time.Duration `mapstructure:"read_timeout"`
			WriteTimeout time.Duration `mapstructure:"write_timeout"`
		} `mapstructure:"http"`
		Database struct {
			DSN string `mapstructure:"dsn"`
		} `mapstructure:"database"`
		Jwt struct {
			Secret             string        `mapstructure:"secret"`
			AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
			RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
		} `mapstructure:"jwt"`
		Logger string `mapstructure:"logger"`
		Router string `mapstructure:"router"`
	} `mapstructure:"auth"`

	Calc struct {
		Orchestrator struct {
			Port                 int `mapstructure:"port"`
			TimeAdditionMs       int `mapstructure:"time_addition_ms"`
			TimeSubtractionMs    int `mapstructure:"time_subtraction_ms"`
			TimeMultiplicationMs int `mapstructure:"time_multiplication_ms"`
			TimeDivisionMs       int `mapstructure:"time_division_ms"`
		} `mapstructure:"orchestrator"`
		Agent struct {
			Workers int           `mapstructure:"workers"`
			Timeout time.Duration `mapstructure:"timeout"`
		} `mapstructure:"agent"`
	} `mapstructure:"calc"`
}
