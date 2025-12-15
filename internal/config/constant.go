package config

import "time"

const (
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
	AuthMiddlewareContextKey   = "claims"
	BcryptCost                 = 3
)

const (
	LoggerUseLocalTime     = false
	LoggerFileMaxSizeInMB  = 100
	LoggerFileMaxAgeInDays = 30
	LoggerMaxBackups       = 5
	LoggerCompress         = true
)