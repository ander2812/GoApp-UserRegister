package main

import (
	//"fmt"

	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {

	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)

	if err != nil {
		panic(err.Error())
	}

	return conexion
}

var temp = template.Must(template.ParseGlob("template/*"))

func main() {
	http.HandleFunc("/", Init)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/information", Information)

	log.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)

}

type Users struct {
	username  string
	firstname string
	lastname  string
	birthdate string
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

	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("SELECT * FROM users")

	if err != nil {
		panic(err.Error())
	}

	user := Users{}
	arregloUsers := []Users{}

	for registros.Next() {

		var username string
		var lastname string
		var firstname string
		var birthdate string
		err = registros.Scan(&username, &lastname, &firstname, &birthdate)
		if err != nil {
			panic(err.Error())
		}
		user.username = username
		user.lastname = lastname
		user.firstname = firstname
		user.birthdate = birthdate

		arregloUsers = append(arregloUsers, user)
	}
	temp.ExecuteTemplate(w, "information", arregloUsers)

	if r.Method == "POST" {

		username := r.FormValue("username")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		birthdate := r.FormValue("birthday")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO users(username,lastname,firstname,birthdate) VALUES('?','?','?','?')")

		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(username, lastname, firstname, birthdate)

		http.Redirect(w, r, "/", 301)

	}

}
