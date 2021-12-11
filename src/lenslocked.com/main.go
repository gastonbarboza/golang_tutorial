package main

import (
	"net/http"
	
	"lenslocked.com/views"
	
	"github.com/gorilla/mux"
)

var homeView *views.View
var contactView *views.View

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w,nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w,nil))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {	
	homeView = views.NewView("bootstrap","views/home.gohtml")
	contactView = views.NewView("bootstrap","views/contact.gohtml")
	
	router := mux.NewRouter()
	
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)

	http.ListenAndServe(":8501", router)
}