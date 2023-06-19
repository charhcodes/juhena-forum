package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// struct for individual posts
type Post struct {
	// ID      int
	Title   string
	Content string
	Time    string
}

// struct for posts
type HomePageData struct {
	Posts []Post
}

// handle registration
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

// handle login + session cookies
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
	var storedPassword []byte // holds the hashed password from the database
	err = db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&storedPassword)
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
	cookie := http.Cookie{
		Name:    "session",
		Value:   sessionID,
		Expires: time.Now().Add(300 * time.Second),
		Path:    "/",
	}

	// Set the cookie in the response headers
	http.SetCookie(w, &cookie)

	// Redirect the user to the homepage
	http.Redirect(w, r, "/", http.StatusFound)
}

// show posts on homepage
func executePosts() ([]Post, error) {
	rows, err := db.Query("SELECT title, content, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.Title, &post.Content, &post.Time)
		if err != nil {
			return nil, err
		}
		// Format the datetime string
		t, err := time.Parse("2006-01-02T15:04:05.999999999-07:00", post.Time)
		if err != nil {
			return nil, err
		}
		post.Time = t.Format("January 2, 2006, 15:04:05")
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// serve homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := executePosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := HomePageData{
		Posts: posts,
	}

	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// serve and create post page
func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve create post page
		http.ServeFile(w, r, "createPost.html")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	titleContent := r.Form.Get("postTitle")
	postContent := r.Form.Get("postContent")

	if titleContent == "" || postContent == "" {
		fmt.Fprintln(w, "Error - please ensure fields aren't empty!")
		return
	}

	// Get selected categories
	categories := r.Form["checkboxes"]

	// Convert categories to a comma-separated string
	categoriesString := strings.Join(categories, ",")
	fmt.Printf("Post categories: %s\n", categoriesString)

	dateCreated := time.Now()

	_, err = db.Exec("INSERT INTO posts (title, content, category_id, created_at) VALUES (?, ?, ?, ?)", titleContent, postContent, categoriesString, dateCreated)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not create post", http.StatusInternalServerError)
		return
	}

	fmt.Println("Post successfully created!")

	// Redirect the user to the homepage
	http.Redirect(w, r, "/", http.StatusFound)
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
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/create-post", createPostHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
