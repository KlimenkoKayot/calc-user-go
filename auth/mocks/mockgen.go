//go:generate mockgen -source=../internal/domain/repository/user_repository.go -destination=repository/user_repository_mock.go -package=repository
//go:generate mockgen -source=../internal/domain/repository/user_repository.go -destination=repository/user_repository_mock.go -package=repository
//go:generate mockgen -source=../internal/domain/ports/token_manager.go -destination=jwt/token_manager_mock.go -package=jwt
//go:generate mockgen -destination=logger/logger_mock.go -package=logger github.com/klimenkokayot/avito-go/libs/logger Logger

package mocks
