package repository

import (
	"log/slog"

	"github.com/ahhhmadtlz/expense-tracker/internal/repository/mysql"
)

type DB struct {
	conn *mysql.MySQLDB
	logger *slog.Logger
}

func New(conn *mysql.MySQLDB,logger *slog.Logger) *DB {
	return &DB{
		conn: conn,
		logger:logger,
	}
}
