package main

import (
	"database/sql"
	f "fmt"
	"html/template"
	h "net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
)

var tpl *template.Template
var db *sql.DB

func main() {
	tpl, _ = template.ParseGlob("templates/*.html")
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/coffeeShopDB")
	if err != nil {
		f.Println("Error validation sql.Open arguments")
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		f.Println("Error verifying connection with db.Ping")
		panic(err.Error())
	}

	h.HandleFunc("/login", loginHandler)
	h.HandleFunc("/loginauth", loginAuthHandler)
	h.HandleFunc("/register", registerHandler)
	h.HandleFunc("/registerauth", registerAuthHandler)
	h.ListenAndServe("localhost:8080", nil)
}

func loginHandler(w h.ResponseWriter, r *h.Request) {
	f.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "/templates/signInPage.html", nil)
}

func loginAuthHandler(w h.ResponseWriter, r *h.Request) {
	f.Println("*****loginAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	f.Printf("username: %v, password: %v", username, password)
	var hash string
	stmt := "select password from users where username = ?"
	row := db.QueryRow(stmt, username)
	err := row.Scan(&hash)
	f.Println("hash from db:", hash)
	if err != nil {
		f.Println("error selecting Hash in db by username")
		tpl.ExecuteTemplate(w, "SignInPage.html", "check username and password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		f.Fprint(w, "You have successfully sign in")
		return
	}
	f.Println("incorrect password")
	tpl.ExecuteTemplate(w, "templates/signInPage.html", "check username and password")
}

func registerHandler(w h.ResponseWriter, r *h.Request) {

}
func registerAuthHandler(w h.ResponseWriter, r *h.Request) {

}
