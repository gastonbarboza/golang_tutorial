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
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	// Create controllers
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)
	
	//Routing code
	router := mux.NewRouter()

	router.Handle("/", staticC.Home).Methods("GET")
	router.Handle("/contact", staticC.Contact).Methods("GET")
	router.HandleFunc("/signup", usersC.New).Methods("GET")
	router.HandleFunc("/signup", usersC.Create).Methods("POST")
	router.Handle("/login", usersC.LoginView).Methods("GET")
	router.HandleFunc("/login",usersC.Login).Methods("POST")

	http.ListenAndServe(fmt.Sprintf(":%d",cfg.Port), router)
}
