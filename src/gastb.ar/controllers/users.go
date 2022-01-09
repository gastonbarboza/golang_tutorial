package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

	"gastb.ar/views"
	"gastb.ar/models"
	"gastb.ar/rand"
)

// UsersController:
// 1. calls the view layer when receiving GET requests on /signup and /login,
// 2. calls the model layer when information is POSTed to /signup and /login. 
type UsersController struct {
	SignupView *views.View
	LoginView  *views.View
	us         *models.UserService
}

// NewUsers creates a new controller on top of an initialized UserService.
func NewUsers(us *models.UserService) *UsersController {
	return &UsersController {
		SignupView: views.NewView("bootstrap", "users/new"),
		LoginView:  views.NewView("bootstrap", "users/login"),
		us:         us,
	}
}

// Form objects and funcitons:

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

// Signup is a handlefunc used to process POST requests on the signup form 
// when a user enters their name, email and password
func (u *UsersController) Signup(w http.ResponseWriter,r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := &models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(user); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/", http.StatusFound)
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
	u.signIn(w, user)
	http.Redirect(w, r, "/", http.StatusFound)
}

// signIn is a method that creates a token, sets it to the user, and
// sets a Cookie header on the ResponseWriter. It returns an error if 
// setting the user was not successful, or generating the token failed
func (u *UsersController) signIn(w http.ResponseWriter, user *models.User) error {
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Token = token
	err = u.us.Update(user)
	if err != nil {
		return err
	}

	cookie := http.Cookie {
		Name:     "remember_token",
		Value:    token,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// CookieTest is used to display cookies set on the current user
func (u *UsersController) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByToken(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	return
	}
	fmt.Fprintln(w, user)
}


