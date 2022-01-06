package models

import "github.com/jinzhu/gorm"

type Stocklist struct {
	gorm.Model
	UserIDs []uint `gorm:"not_null;index"`
	HubID uint `gorm:"not_null"`
}

type stocklistGorm struct {
	db *gorm.DB
}

func (sg *stocklistGorm) Create(stocklist *Stocklist) error {
	return nil
}

type StocklistDB interface {
	Create(stocklist *Stocklist) error
}

type StocklistService interface {
	StocklistDB
}

