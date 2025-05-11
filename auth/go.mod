module github.com/klimenkokayot/avito-go/services/auth

go 1.23.7

require (
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/klimenkokayot/avito-go/libs/jwt v0.0.0-20250429212829-5d8fa29cc7b2
	github.com/klimenkokayot/avito-go/libs/logger v0.0.0-20250429212829-5d8fa29cc7b2
	github.com/klimenkokayot/avito-go/libs/router v0.0.0-20250429212829-5d8fa29cc7b2
	github.com/lib/pq v1.10.9
	go.uber.org/mock v0.5.2
	golang.org/x/crypto v0.37.0
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)
