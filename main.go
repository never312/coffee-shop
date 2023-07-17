package main

import (
	"database/sql"
	f "fmt"
	"html/template"
	h "net/http"
	"unicode"

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

	h.Handle("/static/", h.StripPrefix("/static/", h.FileServer(h.Dir("static"))))

	h.HandleFunc("/", mainPageHandler)
	h.HandleFunc("/login", loginHandler)
	h.HandleFunc("/loginauth", loginAuthHandler)
	h.HandleFunc("/register", registerHandler)
	h.HandleFunc("/registerauth", registerAuthHandler)
	h.ListenAndServe("localhost:8080", nil)
}

func mainPageHandler(w h.ResponseWriter, r *h.Request) {
	f.Println("*****mainPageHandler running*****")
	err := tpl.ExecuteTemplate(w, "mainPage.html", nil)
	if err != nil {
		f.Println("Error executing template:", err)
		return
	}
}

func loginHandler(w h.ResponseWriter, r *h.Request) {
	f.Println("*****loginHandler running*****")
	err := tpl.ExecuteTemplate(w, "signInPage.html", nil)
	if err != nil {
		f.Println("Error executing template:", err)
		return
	}
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
	tpl.ExecuteTemplate(w, "signInPage.html", "check username and password")
}

func registerHandler(w h.ResponseWriter, r *h.Request) {
	f.Println("*****registerHandler running*****")
	err := tpl.ExecuteTemplate(w, "signUpPage.html", nil)
	if err != nil {
		f.Println("Error executing template:", err)
		return
	}
}
func registerAuthHandler(w h.ResponseWriter, r *h.Request) {
	/*
		1. check username criteria
		2. check password criteria
		3. check if username is already exists in database
		4. create bcrypt hash from password
		5. insert username and password hash in database
	*/
	f.Println("*****registerAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	var nameAlphaNumeric = true
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			nameAlphaNumeric = false
		}
	}
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}
	password := r.FormValue("password")
	f.Println("password:", password, "\npswdLength:", len(password))
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			pswdLowercase = true
		case unicode.IsUpper(char):
			pswdUppercase = true
		case unicode.IsNumber(char):
			pswdNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	f.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial, "\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", nameAlphaNumeric, "\nnameLength:", nameLength)
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !nameAlphaNumeric || !nameLength {
		tpl.ExecuteTemplate(w, "signUpPage.html", "please check username and password criteria")
		return
	}
	stmt := "SELECT id FROM users WHERE username = ?"
	row := db.QueryRow(stmt, username)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		f.Println("username already exists, err:", err)
		tpl.ExecuteTemplate(w, "signUpPage.html", "username already taken")
		return
	}
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		f.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "signUpPage.html", "there was a problem registering account")
		return
	}
	f.Println("hash:", hash)
	f.Println("string(hash):", string(hash))
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO users (Username, Hash) VALUES (?, ?);")
	if err != nil {
		f.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "signUpPage.html", "there was a problem registering account")
		return
	}
	defer insertStmt.Close()

	var result sql.Result
	result, err = insertStmt.Exec(username, hash)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	f.Println("rowsAff:", rowsAff)
	f.Println("lastIns:", lastIns)
	f.Println("err:", err)
	if err != nil {
		f.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "signUpPage.html", "there was a problem registering account")
		return
	}
	f.Fprint(w, "congrats, your account has been successfully created")
}
