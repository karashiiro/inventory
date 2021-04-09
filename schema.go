package main

import "gorm.io/gorm"

// PlayerItem represents a player-owned item in some quantity.
type PlayerItem struct {
	gorm.Model
	OwnerID  string `gorm:"index,sort:desc"`
	ItemID   uint32
	Quantity uint32
}
