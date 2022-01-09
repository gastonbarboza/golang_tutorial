package models

import "github.com/jinzhu/gorm"

// stocklistService wraps a database
type StocklistService interface {
	StocklistDB
}

// StocklistDB represents a hub stocklist 
type StocklistDB interface {
	Create(stocklist *Stocklist) error
}

// A stocklist has users that can edit it and an associated hub
type Stocklist struct {
	gorm.Model
	UserIDs uint `gorm:"not_null;index"`
	HubID uint `gorm:"not_null"`
}

// stocklistGorm implements a hub stocklist
type stocklistGorm struct {
	db *gorm.DB
}

// Test to check for correct implementation
var _ StocklistDB = &stocklistGorm{}

func (sg *stocklistGorm) Create(stocklist *Stocklist) error {
	return sg.db.Create(stocklist).Error
}
