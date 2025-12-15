package config

import "github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"

func Default() Config {
	cfx := Config{
		Auth: auth.Config {
			AccessExpirationTime: AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessSubject: AccessTokenSubject,
			RefreshSubject: RefreshTokenSubject,
		},
	}
	return cfx
}