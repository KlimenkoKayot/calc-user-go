package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Logger                 string
	Router                 string
	ServerPort             int
	DatabaseDSN            string
	TokenExpirationMinutes time.Duration
	ReadTimeoutSeconds     time.Duration
	WriteTimeoutSeconds    time.Duration
	JwtSecretKey           string
	AccessTokenExpiration  time.Duration
	RefreshTokenExpiration time.Duration
}

func Load(path string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	readTimeoutString := os.Getenv("READ_TIMEOUT")
	readTimeoutInt, err := strconv.Atoi(readTimeoutString)
	if err != nil {
		return nil, err
	}
	readTimeoutSeconds := time.Second * time.Duration(readTimeoutInt)

	writeTimeoutString := os.Getenv("WRITE_TIMEOUT")
	writeTimeoutInt, err := strconv.Atoi(writeTimeoutString)
	if err != nil {
		return nil, err
	}
	writeTimeoutSeconds := time.Second * time.Duration(writeTimeoutInt)

	databaseDSN := os.Getenv("DATABASE_DSN")

	logger := os.Getenv("LOGGER")
	router := os.Getenv("ROUTER")

	serverPortString := os.Getenv("SERVER_PORT")
	serverPort, err := strconv.Atoi(serverPortString)
	if err != nil {
		return nil, err
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")

	accessTokenExpirationString := os.Getenv("ACCESS_TOKEN_EXPIRATION_TIMEOUT")
	accessTokenExpirationInt, err := strconv.Atoi(accessTokenExpirationString)
	if err != nil {
		return nil, err
	}
	accessTokenExpiration := time.Minute * time.Duration(accessTokenExpirationInt)

	refreshTokenExpirationString := os.Getenv("REFRESH_TOKEN_EXPIRATION_TIMEOUT")
	refreshTokenExpirationInt, err := strconv.Atoi(refreshTokenExpirationString)
	if err != nil {
		return nil, err
	}
	refreshTokenExpiration := time.Hour * time.Duration(refreshTokenExpirationInt)

	return &Config{
		Router:                 router,
		Logger:                 logger,
		ServerPort:             serverPort,
		DatabaseDSN:            databaseDSN,
		ReadTimeoutSeconds:     readTimeoutSeconds,
		WriteTimeoutSeconds:    writeTimeoutSeconds,
		JwtSecretKey:           jwtSecretKey,
		AccessTokenExpiration:  accessTokenExpiration,
		RefreshTokenExpiration: refreshTokenExpiration,
	}, nil
}
