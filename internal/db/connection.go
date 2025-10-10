package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"os"

	"github.com/Shakkuuu/ed-mist-backend/internal/model"
)

type Connection struct {
	DB *gorm.DB
}

// NewConnection データベース接続を作成
func NewConnection(dsn string) (*Connection, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "[GORM] ", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Info,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %w", err)
	}

	// 接続テスト
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("データベース接続テストエラー: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("データベース接続テストエラー: %w", err)
	}

	// マイグレーション実行
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("マイグレーションエラー: %w", err)
	}

	log.Println("データベース接続が確立されました")
	return &Connection{DB: db}, nil
}

// runMigrations マイグレーションを実行
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Device{},
		&model.Stay{},
		&model.Subject{},
		&model.Organization{},
		&model.Room{},
		&model.Lesson{},
	)
}

func (c *Connection) Close() error {
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
