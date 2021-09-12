package main

import (
	//"fmt"
	"log"
	"net/http"
	"text/template"
)

var temp = template.Must(template.ParseGlob("template/*"))

func main() {
	http.HandleFunc("/", Init)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/information", Information)

	log.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)

}

func Init(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Hola develoteca")
	temp.ExecuteTemplate(w, "init", nil)

}

func Create(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Hola develoteca")
	temp.ExecuteTemplate(w, "create", nil)

}

func Information(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Hola develoteca")
	temp.ExecuteTemplate(w, "information", nil)

}
