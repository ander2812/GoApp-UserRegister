package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
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
	http.HandleFunc("/alert", Alert)
	http.HandleFunc("/alert2", Alert2)

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

		if authenticate(r.FormValue("username")) {
			Alert2(w, r)
			http.Redirect(w, r, "/", 301)
		} else {
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

}

//method used to validate if the username is on use
func authenticate(usrname string) bool {
	conectionEnambled := conexionBD()
	storedInfo, err := conectionEnambled.Query("SELECT * FROM users")

	if err != nil {
		panic(err.Error())
	}

	user := User{}
	arrayUser := []User{}

	for storedInfo.Next() {
		var username string
		var password string
		var firstname string
		var lastname string
		var birthdate string
		var country string
		var universidad string

		err = storedInfo.Scan(&username, &password, &firstname, &lastname, &birthdate, &country, &universidad)

		if err != nil {
			panic(err.Error())
		}

		user.Username = username

		arrayUser = append(arrayUser, user)

	}

	for i := 0; i < len(arrayUser); i++ {
		if arrayUser[i].Username == usrname {
			return true
		}
	}

	return false

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

	}

	var status bool

	if r.Method == "POST" {

		for i := 0; i < len(arrayUser) && !status; i++ {

			if arrayUser[i].Username == r.FormValue("usname") {
				if arrayUser[i].Password == r.FormValue("passw") {
					status = true

				}

			}

		}

		if status {
			temp.ExecuteTemplate(w, "information", arrayUser)
		} else {
			Alert(w, r)
			http.Redirect(w, r, "/", 301)
		}

	}

}

func Alert(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "alert", nil)
}

func Alert2(s http.ResponseWriter, t *http.Request) {
	temp.ExecuteTemplate(s, "alert2", nil)
}

func encrypt(txt string, k string) {
	fmt.Println("Encryption Program v0.01")

	text := []byte(txt)
	key := []byte(k)

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		fmt.Println(err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)

	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	fmt.Println(gcm.Seal(nonce, nonce, text, nil))

}

func authenticate2(usrname string, passwr string) bool {
	conectionEnambled := conexionBD()
	storedInfo, err := conectionEnambled.Query("SELECT * FROM users")

	if err != nil {
		panic(err.Error())
	}

	user := User{}
	arrayUser := []User{}

	for storedInfo.Next() {
		var username string

		err = storedInfo.Scan(&username)

		if err != nil {
			panic(err.Error())
		}

		user.Username = username

		arrayUser = append(arrayUser, user)

	}

	for i := 0; i < len(arrayUser); i++ {
		if arrayUser[i].Username == usrname {
			if arrayUser[i].Password == passwr {
				return true
			}
		}
	}

	return false

}
