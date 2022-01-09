package main

import (
	"fmt"
	"net/http"

	"gastb.ar/controllers"
	"gastb.ar/models"
	"gastb.ar/middleware"

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
	userC := controllers.NewUserController(services.UserService)
	requireUserMw := middleware.RequireUser {
		UserService: services.UserService,
	}
	
	// Create intermediate handlers
	profileAuthd := requireUserMw.Apply(staticC.Profile)

	// Routing code
	router := mux.NewRouter()

	router.Handle("/", staticC.Home).Methods("GET")
	router.Handle("/profile", profileAuthd).Methods("GET")
	router.Handle("/signup", userC.SignupView).Methods("GET")
	router.Handle("/login", userC.LoginView).Methods("GET")

	router.HandleFunc("/cookietest",userC.CookieTest).Methods("GET")
	
	router.HandleFunc("/signup", userC.Signup).Methods("POST")
	router.HandleFunc("/login",userC.Login).Methods("POST")

	http.ListenAndServe(fmt.Sprintf(":%d",cfg.Port), router)
}
