package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type IndexPage struct {
	Title string
	Body  string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	page := IndexPage{"Test page", "Hello, World!"}
	tmpl, err := template.ParseFiles("views/layout.html")
	if err != nil {
		// better to redirect to the error page
		panic(err)
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		// better to redirect to the error page
		panic(err)
	}
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
	http.ListenAndServe(":8080", nil)
}
