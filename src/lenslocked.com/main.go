package main

import (
	"net/http"

	"lenslocked.com/controllers"

	"github.com/gorilla/mux"
)

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	router := mux.NewRouter()

	router.Handle("/", staticC.Home).Methods("GET")
	router.Handle("/contact", staticC.Contact).Methods("GET")
	router.HandleFunc("/signup", usersC.New).Methods("GET")
	router.HandleFunc("/signup", usersC.Create).Methods("POST")

	http.ListenAndServe(":8501", router)
}
