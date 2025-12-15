package config

import "github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer `koanf:"http_server"`
	Auth       auth.Config `koanf:"auth"`
}