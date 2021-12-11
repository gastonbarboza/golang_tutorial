package main

import (
	"net/http"
	"html/template"
	
	"github.com/gorilla/mux"
)

var homeTemplate *template.Template
var contactTemplate *template.Template

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func main() {
	var err error
	
	homeTemplate, err = template.ParseFiles("views/home.gohtml","views/layouts/footer.gohtml","views/layouts/header.gohtml")
	if err != nil{
		panic(err)
	}
	
	contactTemplate, err = template.ParseFiles("views/contact.gohtml","views/layouts/footer.gohtml","views/layouts/header.gohtml")
	if err != nil{
		panic(err)
	}
	
	router := mux.NewRouter()
	
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)

	http.ListenAndServe(":8501", router)
}