package forum

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// struct for individual posts
type Post struct {
	id      string
	Title   string
	Content string
	Time    string
	URL     string
}

// struct for posts
type HomePageData struct {
	Posts []Post
}

type PostPageData struct {
	Post *Post
}

var posts []Post

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
