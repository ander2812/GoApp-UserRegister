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
	"fmt"
	"io"
)

func conexionBD() (conexion *sql.DB) {

	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
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
	password  string
	country   string
	college   string
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
		country := r.FormValue("country")
		college := r.FormValue("college")

		conexionEstablecida := conexionBD()

		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO users(username,password,confirmpwd,firstname,lastname,country, college) VALUES(?,?,?,?,?,?)")

		if err != nil {
			panic(err.Error())
		}

		insertarRegistros.Exec(username, password, confirmpwd, firstname, lastname, country, college)

		temp.ExecuteTemplate(w, "init", nil)

		http.Redirect(w, r, "/", 301)

	}

}

func authenticate(usrname string, psw string) bool {
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

		err = storedInfo.Scan(&username, &password)

		if err != nil {
			panic(err.Error())
		}

		user.Username = username
		user.password = password

		arrayUser = append(arrayUser, user)

	}

	for i := 0; i < len(arrayUser); i++ {
		if arrayUser[i].Username == usrname {
			if arrayUser[i].password == psw {
				return true
			}
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

		arrayUser = append(arrayUser, user)

	}
	status := false
	for i := 0; i < len(arrayUser) && !status; i++ {
		if arrayUser[i].Username == r.FormValue("username") {
			if arrayUser[i].password == r.FormValue("password") {
				status = true
			}
		}
	}

	if status {
		temp.ExecuteTemplate(w, "information", arrayUser)
	} else {
		http.Redirect(w, r, "/", 301)
	}

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
