package models

import "github.com/jinzhu/gorm"

// A stocklist will be a list of items associated with a certain user. 
// The user can create more than one.
type Stocklist struct {
	gorm.Model
	UserIDs uint `gorm:"not_null;index"`
}

// StocklistDB represents the database where stocklists are stored.
type StocklistDB interface {
	Create(stocklist *Stocklist) error
}


// stocklistService wraps the database to add other service methods.
type StocklistService struct {
	StocklistDB
}

// stocklistGorm implements a hub stocklist.
// It is only accessible by the StocklistService
type stocklistGorm struct {
	db *gorm.DB
}

// Test to check for correct implementation. Otherwise, code does not compile.
var _ StocklistDB = &stocklistGorm{}

//
// 1. StocklistService methods and related functions
//

//
// 2. StocklistDB methods and related functions
//

// Creates a stocklist and writes it; will only be accessible for logged in
// users, so user information will be available in context.
func (sg *stocklistGorm) Create(stocklist *Stocklist) error {
	return sg.db.Create(stocklist).Error
}
