package main

import (
	"fmt"
	"net/http"

	"lenslocked.com/controllers"
	"lenslocked.com/models"

	"github.com/gorilla/mux"
)

func main() {
	// Config information
	cfg := DefaultConfig()
	psqlInfo := DefaultPostgresConfig().ConnectionInfo()
	
	// Connect to database
	services, err := models.NewServices(psqlInfo)
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
	router.HandleFunc("/signup", usersC.New).Methods("GET")
	router.HandleFunc("/signup", usersC.Create).Methods("POST")
	router.Handle("/login", usersC.LoginView).Methods("GET")
	router.HandleFunc("/login",usersC.Login).Methods("POST")
	router.HandleFunc("/cookietest",usersC.CookieTest).Methods("GET")

	http.ListenAndServe(fmt.Sprintf(":%d",cfg.Port), router)
}
