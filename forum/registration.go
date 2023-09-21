package forum

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	// "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// generate random session ID
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.NewSource(time.Now().UnixNano())
}

func generateUserID(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

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

	if email == "" || username == "" || password == "" {
		http.Error(w, "Please fill out all fields", http.StatusBadRequest)
		return
	}

	// Generate a new session ID for the user
	sessionID := generateUserID(10)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	// Insert the user into the database with the new session ID
	_, err = DB.Exec("INSERT INTO Users (Email, Username, Password, SessionID) VALUES (?, ?, ?, ?)", email, username, hashedPassword, sessionID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}
	fmt.Println("User registered")

	// Set a session cookie to indicate that the user is logged in
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   sessionID,
		Expires: time.Now().Add(1 * time.Hour),
		Path:    "/",
	})

	// Redirect the user to the homepage, where the logout button will be displayed
	http.Redirect(w, r, "/", http.StatusFound)
}

// handle login + session cookies
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user already has an active session
	existingSessionID, _ := r.Cookie("session")
	if existingSessionID != nil {
		// LogoutHandler(w, r)
		// Add your logout logic here if needed
		fmt.Println("User already has an active session")
	}

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "login.html")
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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

	// Check if there is an existing session for this email in the database
	var existingSessionIDDB string
	err = DB.QueryRow("SELECT SessionID FROM Users WHERE Email = ?", email).Scan(&existingSessionIDDB)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existingSessionIDDB != "" {
		// Clear the existing session data from the database
		clearSessionFromDB(existingSessionIDDB)
	}

	var userId int
	var storedPassword []byte // holds the hashed password from the database
	err = DB.QueryRow("SELECT ID, password FROM Users WHERE email = ?", email).Scan(&userId, &storedPassword)
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

	// Generate a new session ID
	newSessionID := generateUserID(10)
	fmt.Printf("loginhandler session id: %v\n", newSessionID)

	// Insert session ID into the database
	_, err = DB.Exec("UPDATE Users SET SessionID = ? WHERE Email = ?", newSessionID, email)
	fmt.Println("Session ID logged")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store the new session ID in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   newSessionID + "&" + strconv.Itoa(userId),
		Expires: time.Now().Add(1 * time.Hour),
		Path:    "/",
	})
	fmt.Println("Session ID stored, user logged in successfully")

	// Redirect the user to the homepage
	http.Redirect(w, r, "/", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the session data from the database
	sessionCookie, err := r.Cookie("session")
	if err == nil {
		sessionID := sessionCookie.Value
		clearSessionFromDB(sessionID)
	}

	// Clear session and user cookies
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1, // Expire immediately
		Path:   "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "user",
		Value:  "",
		MaxAge: -1, // Expire immediately
		Path:   "/",
	})

	fmt.Println("User successfully logged out")
	// Redirect the user to the login page
	http.Redirect(w, r, "/login", http.StatusFound)
}

// clearSessionFromDB removes the session data from the database
func clearSessionFromDB(sessionID string) {
	// Implement code to clear the session from the database by setting SessionID to NULL
	query := "UPDATE Users SET SessionID = NULL WHERE SessionID = ?"
	_, err := DB.Exec(query, sessionID)
	if err != nil {
		// Handle the error, e.g., log it
		log.Println("Error clearing session from the database:", err)
	}
	fmt.Println("Session ID cleared successfully")
}
