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
	// Check session cookie
	checkCookies(w, r)

	// get postID, userID
	postID = getPostID(w, r)
	userID = getUserID(w, r)

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
		err = addLike(userID, postID)
		if err != nil {
			http.Error(w, "Error adding like", http.StatusInternalServerError)
			return
		}
	} else if action == "dislike" {
		err = addDislike(userID, postID)
		if err != nil {
			http.Error(w, "Error adding dislike", http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("Like/Dislike successful!")

	// Redirect back to the post page
	postIDStr := strings.TrimPrefix(r.URL.Path, "/post-like/")
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}

// Other database functions (e.g., addLike, removeLike, checkUserLikeDislike) should be implemented accordingly.

// checkUserLikeDislike checks if a user has already liked or disliked a post.
// func checkUserLikeDislike(userID string, postID int) (liked bool, disliked bool, err error) {
// 	// Execute a query to check if the user has liked the post
// 	row := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = 1", userID, postID)
// 	var likeCount int
// 	if err := row.Scan(&likeCount); err != nil {
// 		return false, false, err
// 	}
// 	liked = likeCount > 0

// 	// Execute a query to check if the user has disliked the post
// 	row = DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = -1", userID, postID)
// 	var dislikeCount int
// 	if err := row.Scan(&dislikeCount); err != nil {
// 		return false, false, err
// 	}
// 	disliked = dislikeCount > 0

// 	return liked, disliked, nil
// }

func checkUserLikeDislike(userID string, postID int) (liked bool, disliked bool, err error) {
	// Execute a query to check if the user has liked or disliked the post
	row := DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = 1", userID, postID)
	var likeCount int
	if err := row.Scan(&likeCount); err != nil {
		return false, false, err
	}
	liked = likeCount > 0

	row = DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = -1", userID, postID)
	var dislikeCount int
	if err := row.Scan(&dislikeCount); err != nil {
		return false, false, err
	}
	disliked = dislikeCount > 0

	// check for data entries that have not been liked or disliked
	row = DB.QueryRow("SELECT COUNT(*) FROM postlikes WHERE user_id = ? AND post_id = ? AND type = 0", userID, postID)
	var neitherCount int
	if err := row.Scan(&neitherCount); err != nil {
		return false, false, err
	}

	// If there is a "neither" entry, treat it as a like and a dislike
	if neitherCount > 0 {
		liked = true
		disliked = true
	}

	return liked, disliked, nil
}

func addLike(userID string, postID int) error {
	existingLike, existingDislike, err := checkUserLikeDislike(userID, postID)
	if err != nil {
		return err
	}

	if existingDislike {
		// Replace value '-1' with '1'
		_, err = DB.Exec("UPDATE postlikes SET type = 1 WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Replaced Dislike with Like")
	} else if existingLike {
		// Toggle value '1' to '0'
		_, err = DB.Exec("UPDATE postlikes SET type = 0 WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Toggled Like to Neither")
	} else {
		// Insert new entry with value '1'
		_, err = DB.Exec("INSERT INTO postlikes (user_id, post_id, type) VALUES (?, ?, 1)", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Added Like")
	}
	return nil
}

func addDislike(userID string, postID int) error {
	existingLike, existingDislike, err := checkUserLikeDislike(userID, postID)
	if err != nil {
		return err
	}

	if existingLike {
		// Replace value '1' with '-1'
		_, err = DB.Exec("UPDATE postlikes SET type = -1 WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Replaced Like with Dislike")
	} else if existingDislike {
		// Toggle value '-1' to '0'
		_, err = DB.Exec("UPDATE postlikes SET type = 0 WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Toggled Dislike to Neither")
	} else {
		// Insert new entry with value '-1'
		_, err = DB.Exec("INSERT INTO postlikes (user_id, post_id, type) VALUES (?, ?, -1)", userID, postID)
		if err != nil {
			return err
		}
		fmt.Println("Added Dislike")
	}
	return nil
}
