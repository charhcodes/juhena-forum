package forum

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var postID int
var userID string

func getPostID(w http.ResponseWriter, r *http.Request) int {
	// Extract post ID from URL
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post-like/")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}
	return postID
}

func getUserID(w http.ResponseWriter, r *http.Request) string {
	// Get user ID from session cookie
	_, userID, err := GetCookieValue(r)
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusBadRequest)
	}
	return userID
}

// Handler for handling like and dislike actions
func HandleLikesDislikes(w http.ResponseWriter, r *http.Request) {
	// get postID, userID
	postID = getPostID(w, r)
	userID = getUserID(w, r)
	fmt.Println(postID)
	fmt.Println(userID)

	r.ParseForm()
	action := r.FormValue("action")

	// Check if the user has already liked/disliked the post
	_, _, err := checkUserLikeDislike(userID, postID)
	if err != nil {
		http.Error(w, "Error checking user like/dislike", http.StatusInternalServerError)
		return
	}

	// Handle the like/dislike action based on user's previous interaction
	// if r.Method == http.MethodPost {
	if action == "like" {
		// remove already existing like/dislike, add like
		err := removeLikeDislike(userID, postID)
		if err != nil {
			http.Error(w, "Error removing like/dislike 1", http.StatusInternalServerError)
			return
		}
		err = addLike(userID, postID)
		if err != nil {
			http.Error(w, "Error adding like 1", http.StatusInternalServerError)
			return
		}
	} else if action == "dislike" {
		// remove already existing like/dislike, add dislikea
		err := removeLikeDislike(userID, postID)
		if err != nil {
			http.Error(w, "Error removing like/dislike 2", http.StatusInternalServerError)
			return
		}
		err = addDislike(userID, postID)
		if err != nil {
			http.Error(w, "Error adding dislike 1", http.StatusInternalServerError)
			return
		}
		// } else {
		// 	err := addLike(userID, postID)
		// 	if err != nil {
		// 		http.Error(w, "Error adding like 2", http.StatusInternalServerError)
		// 		return
		// 	}
	}

	fmt.Println("Like/Dislike successful!")

	// Redirect back to the post page
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post-like/")
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

// Other database functions (e.g., addLike, removeLike, checkUserLikeDislike) should be implemented accordingly.

// checkUserLikeDislike checks if a user has already liked or disliked a post.
func checkUserLikeDislike(userID string, postID int) (liked bool, disliked bool, err error) {
	// Execute a query to check if the user has liked the post
	row := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = 1", userID, postID)
	var likeCount int
	if err := row.Scan(&likeCount); err != nil {
		return false, false, err
	}
	liked = likeCount > 0

	// Execute a query to check if the user has disliked the post
	row = DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = -1", userID, postID)
	var dislikeCount int
	if err := row.Scan(&dislikeCount); err != nil {
		return false, false, err
	}
	disliked = dislikeCount > 0

	return liked, disliked, nil
}

func addLike(userID string, postID int) error {
	_, err := DB.Exec("INSERT INTO postlikes (user_id, post_id, type) VALUES (?, ?, '1')", userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func addDislike(userID string, postID int) error {
	_, err := DB.Exec("INSERT INTO postlikes (user_id, post_id, type) VALUES (?, ?, '-1')", userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func removeLikeDislike(userID string, postID int) error {
	// _, err := DB.Exec("INSERT INTO postlikes (user_id, post_id, type) VALUES (?, ?, 0)", userID, postID)
	_, err := DB.Exec("DELETE FROM postlikes WHERE user_id = ? AND post_id = ?", userID, postID)
	if err != nil {
		return err
	}
	fmt.Println("Like/Dislike removed successfully")
	return nil
}

// func PostLikesHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check session cookie
// 	checkCookies(w, r)

// 	// Get postID from URL path
// 	postIDStr := strings.TrimPrefix(r.URL.Path, "/post-like/")
// 	postID, err := strconv.Atoi(postIDStr)
// 	if err != nil {
// 		http.Error(w, "Invalid post ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Retrieve the sessionID and userID from the cookie
// 	_, userId, err := CommentgetCookieValue(r)
// 	if err != nil {
// 		http.Error(w, "cookie not found", http.StatusBadRequest)
// 		return
// 	}

// 	// insert data into table
// 	_, err = DB.Exec("INSERT INTO postlikes (user_id, post_id) VALUES (?, ?)",
// 		userId, postID)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, "Could not insert data to table", http.StatusInternalServerError)
// 		return
// 	}

// }
