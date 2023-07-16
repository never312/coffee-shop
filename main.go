package main

import (
	"database/sql"
	f "fmt"
	"html/template"
	h "net/http"

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

	insertUser(2, "kikita", "kikita@bk.ru", "12345")
	f.Println("Successful Connecting to Database!")

	h.HandleFunc("/login", loginHandler)
	h.HandleFunc("/loginauth", loginAuthHandler)
	h.HandleFunc("/register", registerHandler)
	h.HandleFunc("/registerauth", registerAuthHandler)
	h.ListenAndServe("localhost:8080", nil)
}

func insertUser(id int, username, email, password string) {
	insert, err := db.Query("INSERT INTO `coffeeShopDB`.`users` (`id`, `username`, `email`, `password`) VALUES (?, ?, ?, ?)", id, username, email, password)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}

func loginHandler(w h.ResponseWriter, r *h.Request) {

}

func loginAuthHandler(w h.ResponseWriter, r *h.Request) {

}
func registerHandler(w h.ResponseWriter, r *h.Request) {

}
func registerAuthHandler(w h.ResponseWriter, r *h.Request) {

}
