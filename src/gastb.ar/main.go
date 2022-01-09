package main

import (
	"fmt"
	"net/http"

	"gastb.ar/controllers"
	"gastb.ar/models"

	"github.com/gorilla/mux"
)

func main() {
	// Config information
	cfg := DefaultConfig()
	psqlInfo := DefaultPostgresConfig().ConnectionInfo()
	hmacSecretKey := cfg.HMAC

	// Connect to database
	services, err := models.NewServices(psqlInfo,hmacSecretKey)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	// Create controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	
	//Routing code
	router := mux.NewRouter()

	router.Handle("/", staticC.Home).Methods("GET")
	router.Handle("/contact", staticC.Contact).Methods("GET")
	router.Handle("/signup", usersC.SignupView).Methods("GET")
	router.Handle("/login", usersC.LoginView).Methods("GET")

	router.HandleFunc("/cookietest",usersC.CookieTest).Methods("GET")
	
	router.HandleFunc("/signup", usersC.Signup).Methods("POST")
	router.HandleFunc("/login",usersC.Login).Methods("POST")

	http.ListenAndServe(fmt.Sprintf(":%d",cfg.Port), router)
}
