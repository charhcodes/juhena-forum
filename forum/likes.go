package forum

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var counter int

// insert like into table when like button is clicked on post
func likePostButton(w http.ResponseWriter, r *http.Request) {
	// if post has already been liked
	if counter == 1 {
		// remove previous like
		_, err := DB.Exec("INSERT INTO postlikes (type) VALUES (0)")
		if err != nil {
			counter = 1
			log.Println(err)
			http.Error(w, "Could not remove like from post", http.StatusInternalServerError)
			return
		} else {
			counter = 0
			return
		}
	}
	// like post
	_, err := DB.Exec("INSERT INTO postlikes (type) VALUES (1)")
	if err != nil {
		counter = 0
		log.Println(err)
		http.Error(w, "Could not like post", http.StatusInternalServerError)
		return
	} else {
		counter = 1
		fmt.Println("Post like successful")
		return
	}
}

// insert dislike into table when dislike button is clicked on post
func dislikePostButton(w http.ResponseWriter, r *http.Request) {
	// if post has already been disliked
	if counter == -1 {
		// remove previous dislike
		_, err := DB.Exec("INSERT INTO postlikes (type) VALUES (0)")
		if err != nil {
			counter = -1
			log.Println(err)
			http.Error(w, "Could not remove dislike from post", http.StatusInternalServerError)
			return
		} else {
			counter = 0
			return
		}
	}
	// dislike post
	_, err := DB.Exec("INSERT INTO postlikes (type) VALUES (-1)")
	if err != nil {
		counter = 0
		log.Println(err)
		http.Error(w, "Could not dislike post", http.StatusInternalServerError)
		return
	} else {
		counter = -1
		fmt.Println("Post dislike successful")
		return
	}
}

func PostLikesHandler(w http.ResponseWriter, r *http.Request) {
	// Check session cookie
	checkCookies(w, r)

	// Get postID from URL path
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Retrieve the sessionID and userID from the cookie
	_, userId, err := CommentgetCookieValue(r)
	if err != nil {
		http.Error(w, "cookie not found", http.StatusBadRequest)
		return
	}

	// insert data into table
	_, err = DB.Exec("INSERT INTO postlikes (user_id, post_id) VALUES (?, ?)",
		userId, postID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not insert data to table", http.StatusInternalServerError)
		return
	}

}
