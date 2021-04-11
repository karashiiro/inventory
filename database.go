package main

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func initDatabase() (*Database, error) {
	dsn := os.Getenv("INVENTORY_PGL_CONNECTION_STRING")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&PlayerItem{})

	return &Database{
		db: db,
	}, nil
}
