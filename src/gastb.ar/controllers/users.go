package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

	"gastb.ar/views"
	"gastb.ar/models"
)

// UsersController :::viewing of signup and login pages and storage and
// retrieval of information from user database. Users 

type UsersController struct {
	SignupView *views.View
	LoginView  *views.View
	us         *models.UserService
}

func NewUsers(us *models.UserService) *UsersController {
	return &UsersController {
		SignupView: views.NewView("bootstrap", "users/new"),
		LoginView:  views.NewView("bootstrap", "users/login"),
		us:         us,
	}
}

type SignupForm struct {
	Name string `schema:"name"`
	Email string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}

// Create is a handlefunc used to process POST requests on the signup form 
// when a user enters their name, email and password
func (u *UsersController) Create(w http.ResponseWriter,r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User is", user)
}

// Login is is a handler used to process POST requests on the login form when
// user sends their email and password
func (u *UsersController) Login(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address.")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password provided.")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	return
	}
	cookie := http.Cookie {
		Name:     "email",
		Value:    user.Email,
		HttpOnly: true,
	}
	http.SetCookie(w,&cookie)
	fmt.Fprintln(w,user)
}

// CookieTest is used to display cookies set on the current user
func (u *UsersController) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Email is:", cookie.Value)
}
