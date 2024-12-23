package main

import (
	"fmt"
	"log"
	"strconv"

	. "go_server/models"
	"go_server/views"

	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	dsn := "postgres://drabart:123456@localhost:5432/go_server_test_db"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Unable to connect: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ping failed: %v\n", err)
	}

	var greeting string
	err = db.QueryRow("SELECT 'Hello, PQSL!'").Scan(&greeting)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}

	fmt.Println(greeting)

	return db
}

func main() {
	courses := []*Course{}
	courses = append(courses, &Course{Name: "CSE", URL: "https://google.com"})

	component := views.Test("Bartosz")
	courses_page := views.CoursesPage(courses)
	// initDB()

	r := mux.NewRouter()

	r.Handle("/", templ.Handler(component))
	r.Handle("/courses", templ.Handler(courses_page))
	r.Handle("/courses/", templ.Handler(courses_page))

	const port = 8080
	fmt.Println("Started server on port", port)
	port_string := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(port_string, r))
}
