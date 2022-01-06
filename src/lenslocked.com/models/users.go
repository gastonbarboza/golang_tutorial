package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
)

// UsersDB is an interface that can interact with the users database.
//
// For single user queries:
// user found returns nil error;
// user not found returns ErrNotFound;
// other errors may also be returned if they arise. 
//
// These "other errors" will result in a 500 error.
type UserDB interface {
	//Query methods
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)

	//Edit methods
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	//Utility methods
	Close() error
	AutoMigrate() error
	DestructiveReset() error
}
// We export the interface so documentation is exported, but we will not 
// export the implementation.

// userGorm is the database interaction layer
// implementing the UserDB interface.
type userGorm struct {
	db *gorm.DB
}

var _ UserDB = &userGorm{}
// Checks to see if userGorm is correctly implemented; otherwise code
// does not compile.

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

// newUserGorm instatiates a userGorm
func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm {
		db: db,
	}, nil
}

// UserService exports the UserDB implementation and implements non-database
// related services.
type UserService struct {
	UserDB
}

func NewUserService(connectionInfo string) (*UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	
	if err != nil {
		return nil, err
	}
	
	return &UserService {
		UserDB: ug,
	}, nil
}

// Close closes the connection to the database.
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// Create will create the provided user and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (ug *userGorm) Create(user *User) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return ug.db.Create(user).Error
}

// We define several errors that may arise user resource manipulation
var (
	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided
	// to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password 
	// is used when attempting to authenticate a user
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

// Auxiliary function that returns first result in database for a query
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// ByID will look up a user with the provided ID.
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be 
// an error generated by the models package.
//
// Any error but ErrNot Found should result in a 500 error.
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail looks up a user with the given email address and
// returns that user.
// Error returns are the same as ByID

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// Authenticate checks validity of email and passowrd
// If the email provided is invalid, it returns 
//   nil, ErrNotFound
// If the password provided is invalid, it returns
//   nil, ErrInvalidPassword
// If all is valid, it returns
//   user, nil
// Otherwise, it returns whatever error arises
//   nil, error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
}

// Update will update the provided user with all of the data
// in the provided user object.
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// AutoMigrate will attempt to automatically migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

//DestructiveReset drops the user table and rebuilds it
func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}
