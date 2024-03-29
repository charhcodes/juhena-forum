package forum

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// struct for individual posts
type Post struct {
	ID         string
	Title      string
	Content    string
	Time       string
	LikesCount int // added JB
	DislikeCount int // added JB
	URL        string
	// Author  string
}

// struct for comments
type Comment struct {
	UserID  string //
	PostID  string //
	Content string
	CommentID string
	Time    string
	LikesCount int // added JB
	DislikeCount int // added JB
	
	// Author  string
}

// struct for posts
type HomePageData struct {
	Posts []Post
}

type PostPageData struct {
	Post *Post
	Comments []Comment  // added
    Success  bool       // For displaying the success message 
}

// struct to contain comments
type CommentsData struct {
	Comment []Comment
}

// var comments []Comment

//var posts []Post

// initialise DB
func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer DB.Close()
}

func Shutdown() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
