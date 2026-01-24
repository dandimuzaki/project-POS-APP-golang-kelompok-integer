package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"project-POS-APP-golang-integer/pkg/utils"

	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(config utils.DatabaseCofig) (*gorm.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
		config.Username, config.Password, config.Name, config.Host)

	// Setup logger for GORM
	newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
					SlowThreshold:             time.Second, // Slow SQL threshold
					LogLevel:                  logger.Info, // Log level
					IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
					Colorful:                  true,        // Disable color
			},
	)

	conn, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: newLogger,
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	sqlDB, err := conn.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return conn, nil
}
