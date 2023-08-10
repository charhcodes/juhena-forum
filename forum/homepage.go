package forum

import (
	"net/http"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func executePosts() ([]Post, error) {
	var posts []Post //local struct - don't change as it duplicates the posts for some reason. 
	rows, err := DB.Query("SELECT id, title, content, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.id, &post.Title, &post.Content, &post.Time)
		if err != nil {
			return nil, err
		}
		// Format the datetime string
		t, err := time.Parse("2006-01-02T15:04:05.999999999-07:00", post.Time)
		if err != nil {
			return nil, err
		}
		post.Time = t.Format("January 2, 2006, 15:04:05")
		// make post URLs
		post.URL = "/post/" + post.id
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// reverse posts
	posts = reverse(posts)
	return posts, nil
}

// reverse posts (latest first)
func reverse(s []Post) []Post {
	//runes := []rune(s)
	length := len(s)
	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// serve homepage
func HomeHandler(w http.ResponseWriter, r *http.Request) {
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
