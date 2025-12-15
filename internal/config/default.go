package config

import (
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
)

func Default() Config {
	cfx := Config{
		Auth: auth.Config {
			AccessExpirationTime: AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessSubject: AccessTokenSubject,
			RefreshSubject: RefreshTokenSubject,
		},
		Logger: logger.Config{
			UseLocalTime:      LoggerUseLocalTime,
			FileMaxSizeInMB:   LoggerFileMaxSizeInMB,
			FileMaxAgeInDays:  LoggerFileMaxAgeInDays,
			MaxBackups:        LoggerMaxBackups,
			Compress:          LoggerCompress,
		},
	}
	return cfx
}