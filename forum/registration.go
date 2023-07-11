package forum

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

// handle registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
	_, err = DB.Exec("INSERT INTO users (Email, Username, Password) VALUES (?, ?, ?)", email, username, hashedPassword)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	fmt.Println("User registered")
	http.Redirect(w, r, "/", http.StatusFound)
}

// handle login + session cookies
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "login.html")
		return
	}

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
		http.Error(w, "Please fill out all fields", http.StatusBadRequest)
		return
	}
	var storedPassword []byte // holds the hashed password from the database
	err = DB.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No user found with this email", http.StatusUnauthorized)
		} else {
			log.Println(err)
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	// Set browser cookies to store login
	http.SetCookie(w, &http.Cookie{
		Name:  "user",
		Value: email,
	})

	// Generate a session ID
	sessionID := uuid.New().String()

	// Create a new cookie with the session ID
	// Create a new cookie with the session ID
	cookie := http.Cookie{
		Name:    "session",
		Value:   sessionID,
		Expires: time.Now().Add(1 * time.Hour), // Extend cookie expiration time
		Path:    "/",
	}

	// Set the cookie in the response headers
	http.SetCookie(w, &cookie)

	// Redirect the user to the homepage
	http.Redirect(w, r, "/", http.StatusFound)
}
