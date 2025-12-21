package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
)

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Port     int    `koanf:"port"`
	Host     string `koanf:"host"`
	DBName   string `koanf:"db_name"`
}

type MySQLDB struct {
	config Config
	db     *sql.DB
}

func (m MySQLDB) Conn() *sql.DB {
	return m.db
}

func (m *MySQLDB) Close() error {
	logger.Info("Closing MySQL database connection")

	if err := m.db.Close(); err != nil {
		logger.Error("Failed to close MySQL database",
			"error", err.Error(),
		)
		return err
	}

	logger.Info("MySQL database connection closed successfully")
	return nil
}

func New(config Config) *MySQLDB {
	logger.Info("Connecting to MySQL database",
		"host", config.Host,
		"port", config.Port,
		"database", config.DBName,
		"username", config.Username,
	)

	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		config.Username, config.Password, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to open MySQL database",
			"error", err.Error(),
			"host", config.Host,
			"port", config.Port,
		)
		panic(fmt.Errorf("can't open mysql db:%v", err))
	}

	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping MySQL database",
			"error", err.Error(),
			"host", config.Host,
			"port", config.Port,
		)
		panic(fmt.Errorf("can't connect to mysql db: %v", err))
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	logger.Info("Successfully connected to MySQL database",
		"host", config.Host,
		"port", config.Port,
		"database", config.DBName,
		"max_open_conns", 10,
		"max_idle_conns", 10,
		"conn_max_lifetime", "3m",
	)

	return &MySQLDB{config: config, db: db}
}