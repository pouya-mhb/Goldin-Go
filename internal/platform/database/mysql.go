package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
)

const defaultConnectionMaxLifetime = 30 * time.Minute

// OpenMySQL opens and verifies a MySQL database connection.
func OpenMySQL(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", MySQLDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("open mysql database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(defaultConnectionMaxLifetime)

	if err := db.PingContext(ctx); err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("ping mysql database: %w; close database: %v", err, closeErr)
		}

		return nil, fmt.Errorf("ping mysql database: %w", err)
	}

	return db, nil
}

// MySQLDSN builds the MySQL data source name from strongly typed configuration.
func MySQLDSN(cfg config.DatabaseConfig) string {
	mysqlCfg := mysql.NewConfig()
	mysqlCfg.User = cfg.User
	mysqlCfg.Passwd = cfg.Password
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	mysqlCfg.DBName = cfg.Name
	mysqlCfg.ParseTime = true
	mysqlCfg.Params = map[string]string{
		"charset":   "utf8mb4",
		"collation": "utf8mb4_unicode_ci",
		"timeout":   "5s",
	}

	return mysqlCfg.FormatDSN()
}
