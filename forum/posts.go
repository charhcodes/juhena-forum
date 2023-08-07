package forum

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	_, userId, err := GetCookieValue(r)
	if err != nil {
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}

	titleContent := r.Form.Get("postTitle")
	postContent := r.Form.Get("postContent")

	if titleContent == "" || postContent == "" {
		fmt.Fprintln(w, "Error - please ensure title and post content fields are not empty!")
		return
	}

	// Get selected categories
	categories := r.Form["checkboxes"]

	// Convert categories to a comma-separated string
	categoriesString := strings.Join(categories, ",")
	fmt.Printf("Post categories: %s\n", categoriesString)

	// dateCreated := time.Now()

	// _, err = DB.Exec("INSERT INTO posts (title, content, category_id, created_at) VALUES (?, ?, ?, ?)", titleContent, postContent, categoriesString, dateCreated)
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "Could not create post", http.StatusInternalServerError)
	// 	return
	// }

	//ricky below
	dateCreated := time.Now()
	fmt.Println(userId)
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Could not convert", http.StatusInternalServerError)
		return
	}
	//ricky below - added placeholders and userintid
	_, err = DB.Exec("INSERT INTO posts (user_id, title, content, category_id, created_at) VALUES (?, ?, ?, ?, ?)", userIdInt, titleContent, postContent, categoriesString, dateCreated)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not create post", http.StatusInternalServerError)
		return
	}

	fmt.Println("Post successfully created!")

	// Redirect the user to the homepage
	http.Redirect(w, r, "/", http.StatusFound)
}

func GetCookieValue(r *http.Request) (string, string, error) {
	//ricky below - indices represent the split to cookie and value
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", "", err
	}
	value := strings.Split(cookie.Value, "&")

	return value[0], value[1], nil
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
	// Get post id from the URL path
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Assuming you have a function to retrieve the logged-in user's ID, if available
	// If not, you can set it to a default value or handle it as you see fit.
	// userID := "user1"

	// Get the post data by calling the getPostByID function or fetching it from the database
	post, err := getPostByID(postIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Assuming your Post struct has a field named PostID
	var data struct {
		PostID   int
		Post     Post
		Comments []Comment
		Success  bool // Add the Success field to indicate if the comment was successfully posted
	}

	data.PostID = postID
	data.Post = *post // Use the dereferenced post pointer
	data.Comments = comments
	data.Success = r.URL.Query().Get("success") == "1"

	tmpl, err := template.ParseFiles("postPage.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template with the data
	err = tmpl.ExecuteTemplate(w, "postPage.html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// func PostPageHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get post id from the URL path
// 	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")

// 	postID, err := strconv.Atoi(postIDStr)
// 	if err != nil {
// 		http.Error(w, "Invalid post ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Execute comments associated with the post
// 	// comments, err := executeComments(postID)
// 	// if err != nil {
// 	// 	http.Error(w, "Could not retrieve comments", http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// Assuming you have a function to fetch a single post by its ID
// 	post, err := getPostByID(postIDStr)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	// Pass the comments data to the template
// 	// data := struct {
// 	// 	Post     Post
// 	// 	Comments []Comment
// 	// }{
// 	// 	Post:     *post, // Use the dereferenced post pointer
// 	// 	Comments: comments,
// 	// }

// 	// Assuming your Post struct has a field named PostID
// 	data := struct {
// 		PostID   int
// 		Post     Post
// 		Comments []Comment
// 	}{
// 		PostID:   postID,
// 		Post:     *post, // Use the dereferenced post pointer
// 		Comments: comments,
// 	}

// 	tmpl, err := template.ParseFiles("postPage.html")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Render the template with the data
// 	err = tmpl.ExecuteTemplate(w, "postPage.html", data) // Use the correct template name
// 	if err != nil {
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

// func PostPageHandler(w http.ResponseWriter, r *http.Request) {
// 	// get post id
// 	postID := strings.TrimPrefix(r.URL.Path, "/post/")

// 	// Assuming you have a function to fetch a single post by its ID
// 	post, err := getPostByID(postID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	data := PostPageData{
// 		Post: post,
// 	}

// 	tmpl, err := template.ParseFiles("postPage.html")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	err = tmpl.Execute(w, data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }
