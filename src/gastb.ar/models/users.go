package models

import (
	"errors"

//	"gastb.ar/rand"
	"gastb.ar/hash"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
)

// User information is coded in a User type and stored in the users database.
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Token        string `gorm:"-"`
	TokenHash    string `gorm:"not null;unique_index"`
}

// UsersDB is an interface that can interact with the users database.
//
// For single user queries:
// user found returns nil error;
// user not found returns ErrNotFound;
// other errors may also be returned if they arise. 
type UserDB interface {
	//Query methods
	ByID(id uint)                 (*User, error)
	ByEmail(email string)         (*User, error)
	ByTokenHash(tokenHash string) (*User, error)

	//Edit methods
	Create(user *User) error
	Update(user *User) error
	Delete(id uint)    error
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

// UserService wraps the UserDB implementation and implements non-database
// related services.
type UserService struct {
	db   UserDB
	hmac hash.HMAC
}

//
// 1. UserService methods and related functions
//

// NewUserService instatiates a UserService on a database connection and
// a hasher for user tokens.
func NewUserService(db *gorm.DB, hmacSecretKey string) *UserService {
	ug := &userGorm{db}
	
	hmac := hash.NewHMAC(hmacSecretKey)

	return &UserService {
		db:     ug,
		hmac:   hmac,
	}
}

// Create takes a user object, hashes sensitive data and passes it on to
// the database layer.
func (us *UserService) Create(user *User) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedBytes)

	if err != nil {
		return err
	}
	return us.db.Create(user)
}

// Update takes a user object, hashes sensitive data and passes it on to
// the database layer
func (us *UserService) Update(user *User) error {
	user.TokenHash = us.hmac.Hash(user.Token)
	return us.db.Update(user)
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
	foundUser, err := us.db.ByEmail(email)
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

// ByToken takes in a token, hashes it, and uses the hash to search 
// for the corresponding user and returns them
// It also returns whatever error is returned by UserDB when searching for 
// that hash.
func (us *UserService) ByToken(token string) (*User, error) {
	tokenHash := us.hmac.Hash(token)
	return us.db.ByTokenHash(tokenHash)
}

//
// 2. UserDB methods and related functions
//

// We define several errors that may arise manipulating users
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

// Create takes a User object and writes it to the database.
// If this results in an error, it returns it.
func (ug *userGorm) Create(user *User) error {
		return ug.db.Create(user).Error
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

// ByID looks up a user with the provided ID and returns them.
// If the user is found, it also returns a nil error.
// If the user is not found, it also returns ErrNotFound.
// If there is another error, it also returns an error with
// more information about what went wrong. This may not be 
// an error generated by the models package.
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

// ByEmail looks up a user with the given email address and returns them.
// Error returns are the same as ByID.
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

// ByToken looks up a user with a given token hash and returns them.
// Error returns are the same as ByID.
func (ug *userGorm) ByTokenHash(tokenHash string) (*User, error) {
	var user User
	db := ug.db.Where("token_hash = ?", tokenHash)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}
