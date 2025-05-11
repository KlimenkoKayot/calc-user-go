package config

import (
	"time"
)

type ApiGatewayConfig struct {
	Http struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"http"`
	Services struct {
		Auth struct {
			URL string `yaml:"url"`
		} `yaml:"auth"`
		Calc struct {
			URL string `yaml:"url"`
		} `yaml:"calc"`
	} `yaml:"services"`
	Router string `yaml:"router"`
	Logger string `yaml:"logger"`
}

type AuthConfig struct {
	Http struct {
		Host         string        `yaml:"host"`
		Port         int           `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
	} `yaml:"http"`
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
	Jwt struct {
		Secret             string        `yaml:"secret"`
		AccessTokenExpiry  time.Duration `yaml:"access_token_expiry"`
		RefreshTokenExpiry time.Duration `yaml:"refresh_token_expiry"`
	} `yaml:"jwt"`
	Logger string `yaml:"logger"`
	Router string `yaml:"router"`
}

type CalcOrchestratorConfig struct {
	Port                 int    `yaml:"port"`
	TimeAdditionMs       uint64 `yaml:"time_addition_ms"`
	TimeSubtractionMs    uint64 `yaml:"time_subtraction_ms"`
	TimeMultiplicationMs uint64 `yaml:"time_multiplication_ms"`
	TimeDivisionMs       uint64 `yaml:"time_division_ms"`
}

type CalcAgentConfig struct {
	Workers int           `yaml:"workers"`
	Timeout time.Duration `yaml:"timeout"`
}

type CalcConfig struct {
	Orchestrator CalcOrchestratorConfig `yaml:"orchestrator"`
	Agent        CalcAgentConfig        `yaml:"agent"`
}

type Config struct {
	ApiGateway ApiGatewayConfig `yaml:"api_gateway"`
	Auth       AuthConfig       `yaml:"auth"`
	Calc       CalcConfig       `yaml:"calc"`
}
