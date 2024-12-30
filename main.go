package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	. "go_server/models"
	"go_server/views"

	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var db *gorm.DB

func initDB() {
	dsn := "postgres://drabart:123456@localhost:5432/go_server_test_db"
	new_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect: %v\n", err)
	}

	// Create tables if missing
	new_db.AutoMigrate(&User{})

	db = new_db
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	loginComponent := views.Login("", "")
	loginComponent.Render(context.Background(), w)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(nickname, password string) bool {
	hash, err := hashPassword(password)
	if err != nil {
		return false
	}

	user := User{}
	result := db.Where("nickname = ?", nickname).First(&user)
	switch result.Error {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		fmt.Println("User did not exist. Creating")
		user.Nickname = nickname
		user.PasswordHash = hash
		user.IsAdmin = false
		db.Create(&user)
	default:
		fmt.Println("Error looking for user")
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

func loginAttemptHandler(w http.ResponseWriter, r *http.Request) {
	nick := r.FormValue("nick")
	password := r.FormValue("password")

	if checkPasswordHash(nick, password) {
		fmt.Println("Logged in successfully")

		session, _ := store.Get(r, "session-login")
		// Set some session values.
		session.Values["logged_in"] = true
		session.Values["username"] = nick

		err := session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		fmt.Println("Wrong password")
		loginComponent := views.Login("Incorrect user or password", nick)
		loginComponent.Render(context.Background(), w)
	}
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	mux.Vars(r)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-login")

	session.Values["logged_in"] = false
	session.Values["nickname"] = nil

	w.Write([]byte("Logged out"))
}

func withSlashHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	http.Redirect(w, r, path, http.StatusMovedPermanently)
}

func checkLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-login")

		if err != nil {
			fmt.Println(err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if loggedIn, ok := session.Values["logged_in"].(bool); !ok || !loggedIn {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func secretHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Secret page!"))
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func main() {
	initDB()

	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	r.HandleFunc("/{path:.+}/", withSlashHandler)
	r.Handle("/", templ.Handler(views.Test("flsdjf"))).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/login", loginAttemptHandler).Methods("POST")
	r.HandleFunc("/logout", logoutHandler).Methods("GET")

	sub := r.NewRoute().Subrouter()
	sub.Use(checkLogin)
	sub.HandleFunc("/secret", secretHandler).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.RequestURI)
		http.NotFound(w, r)
	})

	const port = 8080
	fmt.Println("Started server on port", port)
	port_string := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(port_string, r))
}
