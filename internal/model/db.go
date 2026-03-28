package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"git.woda.ink/Woda_OA/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDB initializes the MySQL database connection using GORM
func InitDB(cfg *config.DatabaseConfig) error {
	// Set GORM log level
	logLevel := logger.Warn
	switch cfg.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix, // "" prefix for all tables
			SingularTable: true,            // use singular table name
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)

	// Ping to verify connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("[DB] MySQL connected successfully")
	return nil
}

// GetDB returns the GORM DB instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
