package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "register.html")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if email == "" || username == "" || password == "" { //checking if any of these fields aee empty
		http.Error(w, "Please fill out all fields - we need to create a form here", http.StatusBadRequest)
		return
	}
	//generates a bcrypt hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	_, err = db.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, hashedPassword)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User registered")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "login.html")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	if email == "" || password == "" {
		http.Error(w, "Please fill out all fields - we need to create a form here", http.StatusBadRequest)
		return
	}
	var storedPassword []byte //holds the hashed password from the database
	err = db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
	if err != nil { //above fetches hashed password from users table in database.
		if err == sql.ErrNoRows {
			http.Error(w, "No user found with this email", http.StatusUnauthorized)
		} else {
			log.Println(err)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	// func comparing hashed password against the database
	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// set browser cookies to store login
	http.SetCookie(w, &http.Cookie{
		Name:  "user",
		Value: email,
	})

	// Generate a session ID
	sessionID := uuid.New().String()

	// Create a new cookie with the session ID
	cookie := http.Cookie{
		Name:    "session",
		Value:   sessionID,
		Expires: time.Now().Add(300 * time.Second), // Set the expiration time for the cookie
		Path:    "/",                               // Set the cookie path to '/' to make it available for all paths
	}

	// Set the cookie in the response headers
	http.SetCookie(w, &cookie)

	// Write a response to indicate the cookie has been set
	fmt.Fprintln(w, "Session cookie set! ")

	fmt.Fprintln(w, "Logged in")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	http.ServeFile(w, r, "home.html")

	// check for correct file path
	fmt.Println("File Path:", filePath)
	//http.FileServer(http.Dir("")).ServeHTTP(w, r)
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Print("fetching...")

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	// http.HandleFunc("/set-session-cookie", setSessionCookieHandler)
	http.HandleFunc("/", formHandler)
	// http.HandleFunc("/login.html", formHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
