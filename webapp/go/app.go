package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

type IndexPage struct {
	Title string
	Body  string
}

type LoginContent struct {
	Message string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: authenticate
	http.Redirect(w, r, "/login", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	content := LoginContent{Message: "Isutterへようこそ!!"}
	tmpl := template.Must(template.ParseFiles("views/layout.html", "views/login.html"))

	err := tmpl.Execute(w, content)
	if err != nil {
		panic(err)
	}
}

func initializeHandler(w http.ResponseWriter, r *http.Request) {
	// impossible to deploy a single binary
	exec.Command(os.Getenv("SHELL"), "-c", "../tools/init.sh").Output()
}

func connect() {
	db, err := sql.Open("mysql", "isucon@/isucon")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/initialize", initializeHandler)
	http.ListenAndServe(":8080", nil)
}
