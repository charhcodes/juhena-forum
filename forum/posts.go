package forum

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// check cookies to see if user is logged in properly
func checkCookies(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		// If there is an error, it means the session cookie was not found
		// Redirect user to login page
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if sessionCookie.Value == "" {
		// If the session cookie is empty, the user is not logged in
		// Redirect user to login page
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
}

// CREATE POSTS FUNCTION
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Check session cookie
	checkCookies(w, r)

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

	_, err = DB.Exec("INSERT INTO posts (title, content, category_id, created_at) VALUES (?, ?, ?, ?)", titleContent, postContent, categoriesString, dateCreated)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not create post", http.StatusInternalServerError)
		return
	}

	fmt.Println("Post successfully created!")

	// Redirect the user to the homepage
	http.Redirect(w, r, "/", http.StatusFound)
}

// get post ID
func getPostByID(postID string) (*Post, error) {
	for i := range posts {
		if posts[i].id == postID {
			return &posts[i], nil
		}
	}
	return nil, errors.New("post not found")
}

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	// get post id
	postID := strings.TrimPrefix(r.URL.Path, "/post/")

	// Assuming you have a function to fetch a single post by its ID
	post, err := getPostByID(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := PostPageData{
		Post: post,
	}

	tmpl, err := template.ParseFiles("postPage.html")
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
