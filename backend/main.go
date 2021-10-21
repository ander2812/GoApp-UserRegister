package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"

	"crypto/aes"
    "crypto/cipher"
    "crypto/rand"
)

func conexionBD() (conexion *sql.DB) {

	Driver := "mysql"
	Usuario := "root"
	Contrasenia := "password"
	Nombre := "system"

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
	http.HandleFunc("/insert", Insert)

	log.Println("Servidor corriendo...")
	http.ListenAndServe(":8080", nil)

}

type User struct {
	Username  string
	Firstname string
	Lastname  string
	Birthdate string
	Country   string
	College   string
}

func Init(w http.ResponseWriter, r *http.Request) {

	temp.ExecuteTemplate(w, "init", nil)

}

func Create(w http.ResponseWriter, r *http.Request) {

	temp.ExecuteTemplate(w, "create", nil)

}

func Insert(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		username := r.FormValue("username")
		password := r.FormValue("password")
		confirmpwd := r.FormValue("confirmpwd")
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		birthdate := r.FormValue("birthdate")

		conexionEstablecida := conexionBD()

		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO users(username,password,confirmpwd,firstname,lastname,birthdate) VALUES(?,?,?,?,?,?)")

		if err != nil {
			panic(err.Error())
		}

		insertarRegistros.Exec(username, password, confirmpwd, firstname, lastname, birthdate)

		temp.ExecuteTemplate(w, "init", nil)

		http.Redirect(w, r, "/", 301)

	}

}

func Information(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()

	registros, err := conexionEstablecida.Query("SELECT * FROM users")

	if err != nil {
		panic(err.Error())
	}

	user := User{}
	arrayUser := []User{}

	for registros.Next() {
		var username string
		var password string
		var confirmpwd string
		var firstname string
		var lastname string
		var birthdate string
		err = registros.Scan(&username, &password, &confirmpwd, &firstname, &lastname, &birthdate)

		if err != nil {
			panic(err.Error())
		}

		user.Username = username
		user.Firstname = firstname
		user.Lastname = lastname
		user.Birthdate = birthdate

		arrayUser = append(arrayUser, user)

	}

	temp.ExecuteTemplate(w, "information", arrayUser)

}

func hashFunc(txt string ) string{

	end := ""
	for i := 0; i < len([]rune(txt)); i++ {
		 
	}
	return ""
}
