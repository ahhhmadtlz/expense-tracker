package config

import (
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/repository/mysql"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer `koanf:"http_server"`
	Mysql      mysql.Config    `koanf:"mysql"`
	Auth       auth.Config `koanf:"auth"`
	Logger     logger.Config `koanf:"logger"`
}