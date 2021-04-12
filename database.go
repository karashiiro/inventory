package main

import (
	"os"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	db *gorm.DB
	c  *cache.Cache
}

func initDatabase() (*Database, error) {
	// Connect to PGL server
	dsn := os.Getenv("INVENTORY_PGL_CONNECTION_STRING")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&PlayerItem{})

	// Connect to Redis server
	r := redis.NewClient(&redis.Options{
		Addr: os.Getenv("INVENTORY_REDIS_LOCATION"),
	})

	c := cache.New(&cache.Options{
		Redis: r,
	})

	return &Database{
		db: db,
		c:  c,
	}, nil
}
