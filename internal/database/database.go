package database

import (
	"fmt"
	"log"
	"os/exec"
	"qna-api/internal/config"
	"qna-api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	err = runMigrations(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка миграций: %w", err)
	}

	log.Println("Миграции базы данных успешно применены")

	return &Database{DB: db}, nil
}

func runMigrations(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)
	
	cmd := exec.Command("goose", "-dir", "internal/database/migrations", "postgres", dsn, "up")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка выполнения миграций: %v\nOutput: %s", err, string(output))
	}
	
	log.Printf("Миграции выполнены: %s", string(output))
	return nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("ошибка получения sql.DB: %w", err)
	}
	return sqlDB.Close()
}