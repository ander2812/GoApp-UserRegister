package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func conexionBD() (conexion *sql.DB) {
	conexion, err := sql.Open("sqlite3", "./base.db")
	conexion.Exec("create table if no exists users(username varchar(40) primary key, password varchar(100), firstname varchar(40), lastname varchar(40), birthdate DATE, country varchar(40), universidad varchar(40))")

	if err != nil {
		panic(err.Error())
	}

	log.Println("base de datos conectada")

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
	Username    string
	Firstname   string
	Password    string
	Confirmpwd  string
	Lastname    string
	Birthdate   string
	Country     string
	Universidad string
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
		firstname := r.FormValue("firstname")
		lastname := r.FormValue("lastname")
		birthdate := r.FormValue("birthdate")
		country := r.FormValue("country")
		universidad := r.FormValue("universidad")

		conexionEstablecida := conexionBD()

		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO users(username,password,firstname,lastname,birthdate,country,universidad) VALUES(?,?,?,?,?,?,?)")

		if err != nil {
			panic(err.Error())
		}

		insertarRegistros.Exec(username, password, firstname, lastname, birthdate, country, universidad)

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
		var firstname string
		var lastname string
		var birthdate string
		var country string
		var universidad string
		err = registros.Scan(&username, &password, &firstname, &lastname, &birthdate, &country, &universidad)

		if err != nil {
			panic(err.Error())
		}

		user.Username = username
		user.Firstname = firstname
		user.Password = password
		user.Lastname = lastname
		user.Birthdate = birthdate
		user.Country = country
		user.Universidad = universidad

		arrayUser = append(arrayUser, user)

		if r.Method == "GET" {

			if user.Username != r.FormValue("usname") || user.Password != r.FormValue("passw") {

				log.Println("usuario incorrecto")

			} else if user.Username == r.FormValue("usname") && user.Password == r.FormValue("passw") {

				log.Println("usuario correcto")

			}
		}
	}

	temp.ExecuteTemplate(w, "information", arrayUser)

}
