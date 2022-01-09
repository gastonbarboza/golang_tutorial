package models

import "github.com/jinzhu/gorm"

type Services struct {
	Stocklist StocklistService
	User      *UserService
	db        *gorm.DB
}

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil { 
		return nil, err
	}
	db.LogMode(true)

	return &Services {
		User:      NewUserService(db),
		Stocklist: &stocklistGorm{},
		db:        db,
	}, nil
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Stocklist{}).Error
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Stocklist{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate() 
}