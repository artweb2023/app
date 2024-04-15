package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	port         = ":3000"
	dbDriverName = "mysql"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	dbx := sqlx.NewDb(db, dbDriverName)

	mux := mux.NewRouter()
	mux.HandleFunc("/", redirectToLogin)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/dashbord", dashboard(dbx))
	mux.HandleFunc("/api/login", authentication(dbx)).Methods(http.MethodPost)
	mux.HandleFunc("/api/logout", logout())
	mux.HandleFunc("/api/change", updateCustomerStatus(dbx))
	fmt.Println("Start server")
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(port, mux)
	log.Println("Request completed successfully")
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(err)
}

func openDB() (*sql.DB, error) {
	return sql.Open(dbDriverName, "user:Querty190@@(localhost:3306)/client?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true")
}
