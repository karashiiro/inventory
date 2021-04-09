package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
}

func initDatabase() (*Database, error) {
	dsn := "host=localhost user=inventory password=inventory dbname=inventory port=9920 sslmode=disable TimeZone=America/Los_Angeles"
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
